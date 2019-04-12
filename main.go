package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/dpb587/go-slack-topic-bot/message"
	"github.com/dpb587/go-slack-topic-bot/message/pairist"
	"github.com/nlopes/slack"
)

func main() {
	logger := log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)
	debug, err := strconv.ParseBool(os.Getenv("SLACK_DEBUG"))
	if err != nil {
		panic(fmt.Errorf("invalid: SLACK_DEBUG: %s", os.Getenv("SLACK_DEBUG")))
	}
	api := slack.New(os.Getenv("SLACK_TOKEN"), slack.OptionLog(logger), slack.OptionDebug(debug))

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	interrupt := message.Join(
		" ",
		message.Conditional(
			pairist.WorkingHours("09:00", "18:00", "America/Los_Angeles"),
			pairist.PeopleInRole{
				Team: "sf-bosh",
				Role: "interrupt",
				People: map[string]string{
					"Luan":    "U0FUK0EBH",
					"Jim":     "U8MRYTRHU",
					"Rebecca": "U8YCN97Q9",
					"Charles": "U03RC8WQ6",
					"Miguel":  "UD46WTA4Q",
					// "Jim":     "U02QZ1E3G",
					"Morgan":  "U04V9L81Y",
					"Belinda": "U5EJ8MQUW",
				},
			},
		),
		message.Conditional(
			pairist.WorkingHours("09:00", "18:00", "Europe/Berlin"),
			message.Literal("<!subteam^S309JAD1P|@openstack-cpi>"),
		),
	)

Loop:
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.MessageEvent:
				logger.Println("MessageEvent")
				if ev.Msg.Type == "message" && ev.Msg.SubType != "message_deleted" && containsHandleText(ev.Msg.Text, "UHNH1AZJM") {
					mention, err := interrupt.Message()
					if err != nil {
						logger.Printf("ERROR: %s\n", err)

						continue
					}

					reply := fmt.Sprintf("^ %s", mention)
					rtm.SendMessage(rtm.NewOutgoingMessage(reply, ev.Msg.Channel, slack.RTMsgOptionTS(ev.Msg.Timestamp)))
				}

			case *slack.RTMError:
				logger.Printf("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				logger.Println("Invalid credentials")
				break Loop

			default:
			}
		}
	}
}

func containsHandleText(text, handle string) bool {
	return strings.Contains(text, fmt.Sprintf("<@%s>", handle))
}
