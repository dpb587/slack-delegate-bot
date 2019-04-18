package topiclookup

import (
	"github.com/dpb587/slack-delegate-bot/interrupt"
	"github.com/dpb587/slack-delegate-bot/message"
	"github.com/nlopes/slack"
	"github.com/pkg/errors"
)

type Interrupt struct {
	API     *slack.Client
	Channel string
}

var _ interrupt.Interrupt = &Interrupt{}

func (i Interrupt) Lookup(m message.Message) ([]interrupt.Interruptible, error) {
	channel := m.InterruptTarget
	if i.Channel != "" {
		channel = i.Channel
	}

	info, err := i.API.GetChannelInfo(channel)
	if err != nil {
		return nil, errors.Wrapf(err, "getting info of channel %s", channel)
	}

	value := info.Topic.Value
	// TODO parse
	// bbl-users // :bbl: *interrupt:* <!subteam^S7E4C41HS|@infrastructureteam> note:* bbl _always_ works
	// bbr // BOSH Backup &amp; Restore | interrupt: <@U08J13EG0> <@UCKK7PZKK> :party_gopher: For PCF/customer specific questions, please ask in the #pcf-backup-restore channel in Pivotal Slack.
	// buildpacks // Interrupt:  `@guillermo` `@ty` `@buildpacks-team` | Lead: `@slevine` | CI: <http://bit.ly/cf-buildpacks|bit.ly/cf-buildpacks> | Java BP: <#C03F5ELTK|java-buildpack> | Hours: 9-6pm EST
	// capi // Can I push: <http://canibump.cfapps.io|canibump.cfapps.io> Interrupt: :whale: <@U0GQNFF8R> <@U056V1DDK> :boom-avocado:  | PM: <@U91NR3Q3T> :spacewhale2: : | Operators are standing by to take your call 9-6 Pacific
	// cf-docs // Questions? Interrupt <@U0JAEKNBH>. Contribute to the Docs! <http://docs.cloudfoundry.org/concepts/contribute.html>
	// cli // Question about Apps or the CC API? Try <#C07C04W4Q|capi> first! Interrupt: <!subteam^S1ZAS8DNY|@cli-team> PM: <@U0CPY3BL2> For contributor discussion, please visit <#CDVP0651P|cli-dev-internal>
	// credhub // Please include your CredHub logs in case of Errors | interrupt: <@U6W2F82B1> <@U8TDZ8VU3> | break glass: `@credhub-team` | PM: <@UDFK4K0KT>, <@UHPMJCXGC>
	return res, nil
}
