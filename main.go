package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dpb587/slack-alias-bot/conditions/hours"
	"github.com/dpb587/slack-alias-bot/interrupt"
	"github.com/dpb587/slack-alias-bot/interrupts"
	"github.com/dpb587/slack-alias-bot/interrupts/conditional"
	"github.com/dpb587/slack-alias-bot/interrupts/pairist"
	"github.com/dpb587/slack-alias-bot/interrupts/union"
	"github.com/dpb587/slack-alias-bot/interrupts/usergroup"
	"github.com/dpb587/slack-alias-bot/message"
	"github.com/nlopes/slack"
)

func main() {
	envDebug := os.Getenv("SLACK_DEBUG")
	debug, err := strconv.ParseBool(envDebug)
	if err != nil && envDebug != "" {
		panic(fmt.Errorf("invalid: SLACK_DEBUG: %s", os.Getenv("SLACK_DEBUG")))
	}

	logger := log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)
	api := slack.New(os.Getenv("SLACK_TOKEN"), slack.OptionLog(logger), slack.OptionDebug(debug))

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	tzLosAngeles, _ := time.LoadLocation("America/Los_Angeles")
	tzBerlin, _ := time.LoadLocation("Europe/Berlin")

	interrupt := union.Interrupt{
		Interrupts: []interrupt.Interrupt{
			conditional.Interrupt{
				When: hours.Condition{
					Location: tzLosAngeles,
					Start:    "09:00",
					End:      "18:00",
				},
				Then: pairist.Interrupt{
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
			},
			conditional.Interrupt{
				When: hours.Condition{
					Location: tzBerlin,
					Start:    "09:00",
					End:      "18:00",
				},
				Then: usergroup.Interrupt{
					ID:    "S309JAD1P",
					Alias: "openstack-cpi",
				},
			},
		},
	}

	{ // cloudfoundry wants a healthcheck port
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello, you've requested: %s\n", r.URL.Path)
		})

		go http.ListenAndServe(":8080", nil)
	}

Loop:
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.MessageEvent:
				logger.Println("MessageEvent")
				if ev.Msg.Type == "message" && ev.Msg.SubType != "message_deleted" && containsHandleText(ev.Msg.Text, "UHNH1AZJM") {
					mentions, err := interrupt.Lookup(message.Message{})
					if err != nil {
						logger.Printf("ERROR: %s\n", err)

						continue
					}

					reply := fmt.Sprintf("^ %s", interrupts.Join(mentions, " "))
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
