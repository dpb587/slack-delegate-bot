# slack-delegate-bot

A bot for pulling others into a conversation.


## Usage

### End-User

Mention the bot and it will start a thread with its configured delegates:

    [#bosh] @dberger: @interrupt what is the answer to life, the universe, and everything?
    [#bosh]  >> @interrupt: ^ @s4heid @langered

If configured, the delegates can be channel-specific:

    [#cf-deployment] @dberger: @interrupt what is the point of 42?
    [#cf-deployment]  >> @interrupt: ^ @cdutra @tv

To pull in the interrupt of another channel, prefix the mention with the channel:

    [#cf-deployment] @dberger: can you deploy Deep Thought?
    [#cf-deployment]  >> @tv: #bosh @interrupt can help
    [#cf-deployment]  >> @interrupt: ^ @s4heid @langered

For private interrupt lookup, direct message (with the channel, if relevant):

    [@interrupt] @dberger: #bosh
    [@interrupt] @interrupt: @s4heid @langered


### Configuration

The bot is configured through one or more YAML configuration files. There must be a top-level `delegatebot` key which contains your configuration and inside it should contain the following keys:

 * `watch` - a list of [conditions](#conditions); if one or more match, it applies
 * `delegate` - a [delegator](#delegators) for who to contact
 * `options`
    * `empty_message` - a custom message to use if there was nobody to delegate to

The following example demonstrates watching in a couple channels and, if mentioned, will respond with delegates from multiple sources and with some conditional behaviors.

```yaml
delegatebot:
  # a list of events for this policy
  watch:
  - target: { channel: C1234567 } # global channel ID
  - target: { channel: C9876543 } # global channel ID

  # defining how to find the users to pull in
  delegate:
    # to support pulling from multiple sources
    union:

    # based on a role from pairist
    - pairist:
        team: bosh-director
        role: interrupt

    # or statically
    - user:
        id: U0FUK0EBH

    # or a static group
    - usergroup:
        id: S309JAD1P
        alias: openstack-cpi

    # or only during business hours
    - if:
        when:
        - hours: { tz: America/Los_Angeles, start: 09:00, end: 18:00 }
        - day: { tz: America/Los_Angeles, days: [ Mon, Tue, Wed, Thu, Fri ] }
        then:
          pairist:
              team: bosh-director
              role: interrupt
```


#### Delegators

Delegators are configured to lookup someone who can be contacted.


##### `topiclookup`

Refer to a channel's topic to try and identify the interrupts. The following conventions are matched:

 * `interrupt: <@1>...`
 * `interrupt <@n>...`

```yaml
topiclookup:
  channel: C12345678
```


##### `coalesce`

List multiple interrupts and the first one which finds an interrupt will be used.

```yaml
coalesce:
- when:
  - day: { days: Mon, Wed, Fri }
  then:
    user: { id: U12345678 }
- when:
  - day: { days: Tue, Thu }
  then:
    user: { id: U98765432 }
```


##### `conditional`

Wrap an interrupt in conditional behavior. When multiple conditionals are configured, all must evaluate to true. The `else` behavior is optional.

```yaml
if:
  when:
  - hours: { start: 08:00, end: 18:00 }
  then:
    user: { id: U12345678 }
  else:
    literal: { text: "Try pinging us during work hours" }
```


##### `literal`

Instead of a user or user group, mention an interrupt with literal text.

```yaml
literal: { text: "find the person with the *ninja* hat" }
```


##### `literalmap`

Convert literal interrupts generated from another interrupt source into Slack users or usergroups.

```yaml
literalmap:
  from: { pairist: { team: bosh-director, role: Interrupt } }
  users:
    Danny: U0FUK0EBH
```


##### `pagerduty`

Refer to a PagerDuty escalation policy to find current on-call users. By default, only the first escalation level is used.

```yaml
pagerduty:
  api_key: # literal or $PAGERDUTY_TEAM_x
  escalation_policy: PZI9P8E
  # escalation_level: 0 # to show all users, or a specific level number
```


##### `pairist`

Refer to a team's [pairist](https://pair.ist/) to find people with a particular role.

```yaml
pairist:
  team: bosh-director
  # password: literal # OR $PAIRIST_TEAM_x
  role: Interrupt
  # track: Community
```


##### `union`

List multiple interrupts and all discovered interrupts will be suggested.

```yaml
union:
- user: { id: U12345678 }
- usergroup: { id: S23456789, alias: "slackgroupname" }
```


##### `user`

Interrupt a specific user.

```yaml
user: { id: U12345678 }
```


##### `usergroup`

Interrupt a specific user group.

```yaml
usergroup: { id: S12345678, alias: "slackgroupname" }
```


#### Conditions

Conditions are used to require context-specific details for something can occur.


##### `and`

A list of conditionals, all of which must be satisfied.

```yaml
and:
- day: { days: [ Mon, Tue, Wed, Thu, Fri ] }
- hours: { start: 09:00, end: 18:00 }
```


##### `date`

Satisfied when the current date matches one of the listed dates (`YYYY-MM-DD` format).

```yaml
date:
  tz: America/Los_Angeles
  dates: [ 2019-01-01, 2019-01-21, 2019-02-18, 2019-05-27 ]
```


##### `day`

Satisfied when the current day matches one of the listed days (`Mon` format).

```yaml
day:
  tz: Europe/Berlin
  days: [ Mon, Tue, Wed, Thu, Fri ]
```


##### `hours`

Satisfied when the current time is within a start and end time (`HH:MM` 24h format).

```yaml
hours:
  tz: America/Toronto
  start: 09:00
  end: 17:00
```


##### `not`

Inverts another conditional's result.

```yaml
not:
  day: { days: [ Mon, Tue, Wed, Thu, Fri ] }
```


##### `target`

Satisfied when targeting a specific channel.

```yaml
target:
  channel: C02HPPYQ2
```


##### `or`

A list of conditionals, one of which must be satisfied.

```yaml
or:
- hours: { start: 09:00, end: 12:30 }
- hours: { start: 13:30, end: 17:00 }
```



### Environment Preparation

Before running the app, you may want to create a [slack bot user](https://api.slack.com/bot-users) if you haven't already and invite the bot to a channel. Additionally, the app expects that the following environment variable is present:

 * `SLACK_TOKEN` - an API token for accessing [Slack](https://slack.com) (create one [here](https://apps.slack.com/apps/A0F7YS25R-bots))



### Execution

Run the bot...

```bash
$ go run ./cmd/delegatebot \
    --config=app/cloudfoundry/config/*.yml \
    --config=app/cloudfoundry/config/default.delegatebot run
```


### Deployment

See [`app/cloudfoundry`](app/cloudfoundry) for an example.


## License

[MIT License](./LICENSE)
