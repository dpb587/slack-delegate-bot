package cmd

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/dpb587/slack-delegate-bot/cmd/delegatebot/opts"
	"github.com/nlopes/slack"
)

type StatsCmd struct {
	*opts.Root `no-flags:"true"`

	BotID string `long:"bot-id" description:"Bot ID"`

	userMap map[string]string
}

func (c *StatsCmd) getUserName(api *slack.Client, id string) string {
	if c.userMap == nil {
		c.userMap = map[string]string{}
	}

	if _, found := c.userMap[id]; !found {
		info, err := api.GetUserInfo(id)
		if err != nil {
			if err.Error() == "user_not_found" {
				info = &slack.User{Name: "unknown"}
			} else {
				panic(err)
			}
		}

		c.userMap[id] = info.Name
	}

	return c.userMap[id]
}

func (c *StatsCmd) Execute(_ []string) error {
	slackAPI := c.Root.GetSlackAPI()

	bot, err := slackAPI.GetUserInfo(c.BotID)
	if err != nil {
		panic(err)
	}

	since, _ := time.Parse(time.RFC3339, "2019-05-01T00:00:00Z")
	until, _ := time.Parse(time.RFC3339, "2019-07-06T23:59:59Z")

	searchParams := slack.SearchParameters{
		Sort:          "timestamp",
		SortDirection: "desc",
		Count:         100,
	}

	var finished bool

	for {
		msgs, _, err := slackAPI.Search(fmt.Sprintf("@%s", bot.ID), searchParams)
		if err != nil {
			panic(err)
		}

		for _, msg := range msgs.Matches {
			timestamp := convertSlackTimestamp(msg.Timestamp)

			if timestamp.After(until) {
				continue
			} else if timestamp.Before(since) {
				finished = true

				break
			} else if !strings.Contains(msg.Text, fmt.Sprintf("<@%s>", bot.ID)) {
				// avoid cross-posted/shared messages referencing a mention from elsewhere
				continue
			}

			if msg.Channel.Name != "credhub" || msg.Username != "tviehman" {
				//continue
			}

			event := statEvent{
				Timestamp:   timestamp,
				ChannelID:   msg.Channel.ID,
				ChannelName: msg.Channel.Name,
				UserID:      msg.User,
				UserName:    msg.Username,
				Permalink:   msg.Permalink,
			}

			repliesParams := &slack.GetConversationRepliesParameters{
				ChannelID: msg.Channel.ID,
				Timestamp: msg.Timestamp,
				Limit:     100,
			}

			tsSplit := strings.SplitN(msg.Permalink, "thread_ts=", 2)
			if len(tsSplit) == 2 {
				// search matches may be from within threads and the api doesn't give a separate filed :(
				repliesParams.Timestamp = tsSplit[1]
			}

			for {
				replies, _, nextCursor, err := slackAPI.GetConversationReplies(repliesParams)
				if err != nil {
					panic(err)
				}

				for _, reply := range replies {
					replyTimestamp := convertSlackTimestamp(reply.Timestamp)

					if replyTimestamp.Before(timestamp) || replyTimestamp.Equal(timestamp) {
						// wait until replies are after are initiator
						continue
					}

					if reply.User == c.BotID {
						// always ignore bot responses
						continue
					}

					username := c.getUserName(slackAPI, reply.User)
					delay := replyTimestamp.Sub(event.Timestamp).Seconds()

					event.LastReplyExists = true
					event.LastReplyUserID = &reply.User
					event.LastReplyUserName = &username
					event.LastReplyTimestamp = &replyTimestamp
					event.LastReplyDelay = &delay

					if !event.FirstReplyExists {
						if reply.User == event.UserID {
							// ignore initial asker
							continue
						}

						event.FirstReplyExists = true
						event.FirstReplyUserID = event.LastReplyUserID
						event.FirstReplyUserName = event.LastReplyUserName
						event.FirstReplyTimestamp = event.LastReplyTimestamp
						event.FirstReplyDelay = event.LastReplyDelay
					}
				}

				if nextCursor == "" {
					break
				}

				repliesParams.Cursor = nextCursor

				// lazy throttling avoidance
				time.Sleep(5 * time.Second)
			}

			eventJSON, err := json.Marshal(event)
			if err != nil {
				panic(err)
			}

			fmt.Printf("%s\n", eventJSON)

			// lazy throttling avoidance
			time.Sleep(time.Second)
		}

		if finished {
			break
		}

		searchParams.Page = searchParams.Page + 1

		// lazy throttling avoidance
		time.Sleep(15 * time.Second)
	}

	return nil
}

type statEvent struct {
	Timestamp           time.Time  `json:"timestamp"`
	ChannelID           string     `json:"channel_id"`
	ChannelName         string     `json:"channel_name"`
	UserID              string     `json:"user_id"`
	UserName            string     `json:"user_name"`
	Permalink           string     `json:"permalink"`
	FirstReplyExists    bool       `json:"first_reply_exists"`
	FirstReplyUserID    *string    `json:"first_reply_user_id"`
	FirstReplyUserName  *string    `json:"first_reply_user_name"`
	FirstReplyTimestamp *time.Time `json:"first_reply_timestamp"`
	FirstReplyDelay     *float64   `json:"first_reply_delay"`
	LastReplyExists     bool       `json:"last_reply_exists"`
	LastReplyUserID     *string    `json:"last_reply_user_id"`
	LastReplyUserName   *string    `json:"last_reply_user_name"`
	LastReplyTimestamp  *time.Time `json:"last_reply_timestamp"`
	LastReplyDelay      *float64   `json:"last_reply_delay"`
}

func convertSlackTimestamp(timestamp string) time.Time {
	timeFloat, err := strconv.ParseFloat(timestamp, 10)
	if err != nil {
		panic(err)
	}

	sec, dec := math.Modf(timeFloat)

	return time.Unix(int64(sec), int64(dec*(1e9)))
}
