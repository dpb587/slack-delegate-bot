# Delegators

Delegators are configured to lookup someone who can be contacted.


## `topiclookup`

Refer to a channel's topic to try and identify the interrupts. If a channel is not configured, the channel from the message will be used. The following conventions are matched:

 * `interrupt: <@1>...`
 * `interrupt <@n>...`

```yaml
topiclookup:
  channel: C12345678
```


## `coalesce`

List multiple interrupts and the first one which finds an interrupt will be used.

```yaml
coalesce:
- if:
    when:
    - day: { days: Mon, Wed, Fri }
    then:
      user: { id: U12345678 }
- if:
    when:
    - day: { days: Tue, Thu }
    then:
      user: { id: U98765432 }
```


## `conditional`

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


## `emaillookupmap`

Attempt to map email addresses to Slack users.

```yaml
emaillookupmap:
  from: { pagerduty: { api_key: $PAGERDUTY_API_KEY, escalation_policy: PZI9P8E } }
```


## `literal`

Instead of a user or user group, mention an interrupt with literal text.

```yaml
literal: { text: "find the person with the *ninja* hat" }
```


## `literalmap`

Convert literal interrupts generated from another interrupt source into Slack users or usergroups.

```yaml
literalmap:
  from: { pairist: { team: bosh-director, role: Interrupt } }
  users:
    Danny: U0FUK0EBH
```


## `pagerduty`

Refer to a PagerDuty escalation policy to find current on-call users. By default, only the first escalation level is used.

```yaml
pagerduty:
  api_key: # literal or $PAGERDUTY_TEAM_x
  escalation_policy: PZI9P8E
  # escalation_level: 0 # to show all users, or a specific level number
```


## `pairist`

Refer to a team's [pairist](https://pair.ist/) to find people with a particular role.

```yaml
pairist:
  team: bosh-director
  # password: literal # OR $PAIRIST_TEAM_x
  role: Interrupt
  # track: Community
```


## `union`

List multiple interrupts and all discovered interrupts will be suggested.

```yaml
union:
- user: { id: U12345678 }
- usergroup: { id: S23456789, alias: "slackgroupname" }
```


## `user`

Interrupt a specific user.

```yaml
user: { id: U12345678 }
```


## `usergroup`

Interrupt a specific user group.

```yaml
usergroup: { id: S12345678, alias: "slackgroupname" }
```
