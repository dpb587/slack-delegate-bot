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

The bot is configured through one or more configuration files which look like the following:

```yaml
# a list of conditions for this interrupt policy
when:
- target: { channel: bosh } # global channel ID

# a list of users to interrupt
then:

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


#### Interrupts

Interrupts are used to lookup someone who can be contacted.


##### `channeltopic`

Refer to a channel's topic to try and identify the interrupts. The following conventions are matched:

 * `interrupt: <@1>...`
 * `interrupt <@n>...`

```yaml
channeltopic:
  id: C12345678
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


##### `pairist`

Refer to a team's [pairist](https://pair.ist/) to find people with a particular role.

```yaml
pairist:
  team: bosh-director
  role: Interrupt
  people:
    Charles: U03RC8WQ6
```


##### `union`

List multiple interrupts and all discovered interrupts will be suggested.

```yaml
union:
- user: { id: U12345678 }
- usergroup: { id: S23456789 }
```


##### `user`

Interrupt a specific user.

```yaml
user: { id: U12345678 }
```


##### `usergroup`

Interrupt a specific user group.

```yaml
usergroup: { id: S12345678 }
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


### Execution

Run the bot...

```
$ go run ./main --handler=example/cloudfoundry/*.yml --handler=example/cloudfoundry/default/global.yml run
```


## License

[MIT License](./LICENSE)
