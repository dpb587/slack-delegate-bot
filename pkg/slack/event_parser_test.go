package slack_test

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"

	"github.com/dpb587/slack-delegate-bot/pkg/message"
	. "github.com/dpb587/slack-delegate-bot/pkg/slack"
	"github.com/dpb587/slack-delegate-bot/pkg/slack/slackfakes"
)

var _ = Describe("EventParser", func() {
	const appID = "A1234567"
	const teamID = "T1234567"
	const botUserID = "U1234567"
	const realUserID = "U9876543"

	const localChannelID = "C1234567"
	const remoteChannelID = "C9876543"

	var subject *EventParser
	var fakeUserLookupSlackAPI *slackfakes.FakeUserLookupSlackAPI
	var eventRaw slackevents.EventsAPIEvent

	BeforeEach(func() {
		fakeUserLookupSlackAPI = &slackfakes.FakeUserLookupSlackAPI{}
		fakeUserLookupSlackAPI.GetUserInfoStub = func(in string) (*slack.User, error) {
			if in == botUserID {
				return &slack.User{
					Profile: slack.UserProfile{
						ApiAppID: appID,
					},
				}, nil
			}

			return &slack.User{}, nil
		}

		subject = NewEventParser(NewUserLookup(fakeUserLookupSlackAPI))
		eventRaw = slackevents.EventsAPIEvent{
			APIAppID: appID,
			TeamID:   teamID,
		}
	})

	Describe("ParseAppMention", func() {
		var event slackevents.AppMentionEvent

		BeforeEach(func() {
			event = slackevents.AppMentionEvent{
				User:            realUserID,
				UserTeam:        teamID,
				Channel:         localChannelID,
				Text:            fmt.Sprintf("hi <@%s> i haz questions", realUserID),
				TimeStamp:       "1588524033.0",
				ThreadTimeStamp: "1588524033.1",
			}
		})

		It("parses a default test message", func() {
			msg, reply, err := subject.ParseAppMention(eventRaw, event)
			Expect(err).ToNot(HaveOccurred())
			Expect(reply).To(BeTrue())
			Expect(msg.ChannelTeamID).To(Equal(teamID))
			Expect(msg.ChannelID).To(Equal(localChannelID))
			Expect(msg.UserTeamID).To(Equal(teamID))
			Expect(msg.UserID).To(Equal(realUserID))
			Expect(msg.TargetChannelTeamID).To(Equal(teamID))
			Expect(msg.TargetChannelID).To(Equal(localChannelID))
			Expect(msg.RawTimestamp).To(Equal("1588524033.0"))
			Expect(msg.RawThreadTimestamp).To(Equal("1588524033.1"))
			Expect(msg.RawText).To(Equal(fmt.Sprintf("hi <@%s> i haz questions", realUserID)))
			Expect(msg.Type).To(Equal(message.ChannelMessageType))
			Expect(msg.Time.Format(time.RFC3339)).To(Equal("2020-05-03T16:40:33Z"))
		})

		It("ignores messages from self", func() {
			event.User = botUserID

			_, reply, err := subject.ParseAppMention(eventRaw, event)
			Expect(err).ToNot(HaveOccurred())
			Expect(reply).To(BeFalse())
		})

		It("captures contextual channels", func() {
			event.Text = fmt.Sprintf("hi <#%s> <@%s> i haz questions", remoteChannelID, botUserID)

			msg, reply, err := subject.ParseAppMention(eventRaw, event)
			Expect(err).ToNot(HaveOccurred())
			Expect(reply).To(BeTrue())
			Expect(msg.ChannelID).To(Equal(localChannelID))
			Expect(msg.TargetChannelID).To(Equal(remoteChannelID))
		})

		PIt("parses attachment fallback", func() {
			// msg.Attachments = []slack.Attachment{
			// 	{
			// 		Fallback: msg.Text,
			// 	},
			// }

			// msg.Text = "something else entirely"

			// res, err := subject.ParseMessage(msg)
			// Expect(err).NotTo(HaveOccurred())
			// Expect(res).ToNot(BeNil())
			// Expect(res.Text).To(Equal("something else entirely\n\n---\n\nhelp me, <@U1234567> you're my only hope."))
		})
	})

	Describe("ParseMessage", func() {
		const directID = "D1234567"

		var event slackevents.MessageEvent

		BeforeEach(func() {
			event = slackevents.MessageEvent{
				User:            realUserID,
				UserTeam:        teamID,
				Channel:         directID,
				Text:            fmt.Sprintf("hi <#%s> <@%s> i haz questions", remoteChannelID, botUserID),
				TimeStamp:       "1588524033.0",
				ThreadTimeStamp: "1588524033.1",
				ChannelType:     "mpim",
			}
		})

		It("ignores messages from self", func() {
			event.User = botUserID

			_, reply, err := subject.ParseMessage(eventRaw, event)
			Expect(err).ToNot(HaveOccurred())
			Expect(reply).To(BeFalse())
		})

		It("ignore mentions without a channel", func() {
			event.Text = fmt.Sprintf("hi <@%s> i haz questions", botUserID)

			_, reply, err := subject.ParseMessage(eventRaw, event)
			Expect(err).ToNot(HaveOccurred())
			Expect(reply).To(BeFalse())
		})

		It("ignore channels without a mention", func() {
			event.Text = fmt.Sprintf("tell me about <#%s> please", remoteChannelID)

			_, reply, err := subject.ParseMessage(eventRaw, event)
			Expect(err).ToNot(HaveOccurred())
			Expect(reply).To(BeFalse())
		})

		It("ignore channel mention syntax with non-bot user", func() {
			event.Text = fmt.Sprintf("hi <#%s> <@%s> i haz questions", remoteChannelID, realUserID)

			_, reply, err := subject.ParseMessage(eventRaw, event)
			Expect(err).ToNot(HaveOccurred())
			Expect(reply).To(BeFalse())
		})

		PIt("parses attachment fallback", func() {
			// msg.Attachments = []slack.Attachment{
			// 	{
			// 		Fallback: msg.Text,
			// 	},
			// }

			// msg.Text = "something else entirely"

			// res, err := subject.ParseMessage(msg)
			// Expect(err).NotTo(HaveOccurred())
			// Expect(res).ToNot(BeNil())
			// Expect(res.Text).To(Equal("something else entirely\n\n---\n\nhelp me, <@U1234567> you're my only hope."))
		})

		Context("overhearing mentions", func() {
			It("parses a default test message", func() {
				msg, reply, err := subject.ParseMessage(eventRaw, event)
				Expect(err).ToNot(HaveOccurred())
				Expect(reply).To(BeTrue())
				Expect(msg.ChannelTeamID).To(Equal(teamID))
				Expect(msg.ChannelID).To(Equal(directID))
				Expect(msg.UserTeamID).To(Equal(teamID))
				Expect(msg.UserID).To(Equal(realUserID))
				Expect(msg.TargetChannelTeamID).To(Equal(teamID))
				Expect(msg.TargetChannelID).To(Equal(remoteChannelID))
				Expect(msg.RawTimestamp).To(Equal("1588524033.0"))
				Expect(msg.RawThreadTimestamp).To(Equal("1588524033.1"))
				Expect(msg.RawText).To(Equal(fmt.Sprintf("hi <#%s> <@%s> i haz questions", remoteChannelID, botUserID)))
				Expect(msg.Type).To(Equal(message.ChannelMessageType))
				Expect(msg.Time.Format(time.RFC3339)).To(Equal("2020-05-03T16:40:33Z"))
			})

			It("captures contextual channels", func() {
				event.Text = fmt.Sprintf("hi <#%s> <@%s> i haz questions", remoteChannelID, botUserID)

				msg, reply, err := subject.ParseMessage(eventRaw, event)
				Expect(err).ToNot(HaveOccurred())
				Expect(reply).To(BeTrue())
				Expect(msg.ChannelID).To(Equal(directID))
				Expect(msg.TargetChannelID).To(Equal(remoteChannelID))
			})
		})

		Context("direct", func() {
			BeforeEach(func() {
				event.ChannelType = "im"
				event.Text = fmt.Sprintf("tell me about <#%s> please", remoteChannelID)
			})

			It("parses a default test message", func() {
				msg, reply, err := subject.ParseMessage(eventRaw, event)
				Expect(err).ToNot(HaveOccurred())
				Expect(reply).To(BeTrue())
				Expect(msg.ChannelTeamID).To(Equal(teamID))
				Expect(msg.ChannelID).To(Equal(directID))
				Expect(msg.UserTeamID).To(Equal(teamID))
				Expect(msg.UserID).To(Equal(realUserID))
				Expect(msg.TargetChannelTeamID).To(Equal(teamID))
				Expect(msg.TargetChannelID).To(Equal(remoteChannelID))
				Expect(msg.RawTimestamp).To(Equal("1588524033.0"))
				Expect(msg.RawThreadTimestamp).To(Equal("1588524033.1"))
				Expect(msg.RawText).To(Equal(fmt.Sprintf("tell me about <#%s> please", remoteChannelID)))
				Expect(msg.Type).To(Equal(message.DirectMessageMessageType))
				Expect(msg.Time.Format(time.RFC3339)).To(Equal("2020-05-03T16:40:33Z"))
			})

			It("ignores messages without a channel", func() {
				event.Text = fmt.Sprintf("hi <@%s> i haz questions", realUserID)

				_, reply, err := subject.ParseMessage(eventRaw, event)
				Expect(err).ToNot(HaveOccurred())
				Expect(reply).To(BeFalse())
			})
		})
	})
})
