# slack-alias-bot

When mentioning `bosh-interrupt` the go-alias-bot starts a thread and mentions the interrupt pair.


## Usage

### End-User

If the current channel has configured an interrupt, the bot can be mentioned directly:

    [#bosh] @dberger: @interrupt what is the answer to life, the universe, and everything?
    [#bosh]  >> @interrupt: ^ @s4heid @langered

To pull in the interrupt of another channel, suffix the mention with the channel:

    [#bosh] @dberger: @interrupt #cf-deployment what is the point of 42?
    [#bosh]  >> @interrupt: ^ @cdutra @tv

To ask for an interrupt, direct message with the target:

    [@interrupt] @dberger: #bosh
    [@interrupt] @interrupt: @s4heid @langered


### Configuration

```yaml
# a list of conditions for this interrupt policy
when:
- target: { channel: bosh } # global channel ID

# a list of users to interrupt
interrupt:

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

### Execution

Run the bot...

```
$ slack-alias-bot run config/
```

Validate configuration...

```
$ slack-alias-bot validate config/
```

Environment:

 * **`SLACK_TOKEN`** -- API token
 * `SLACK_DEBUG` -- optionally enable extra debug logging by setting to `true`


## License

[MIT](./LICENSE)
