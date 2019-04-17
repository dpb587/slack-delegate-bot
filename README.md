# slack-alias-bot

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
$ go run ./main --handler=example/yaml/* run
```


## License

[MIT](./LICENSE)
