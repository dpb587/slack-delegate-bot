# slack-delegate-bot

A conventional bot for pulling others into a conversation.

## User Experience

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

## Docs

 * [Slack Setup](docs/slack-setup.md)
 * [Deployment](docs/deployment.md)
 * Configuration
    * [YAML Files](docs/handlers/yaml-config.md)
    * [Conditions](docs/handlers/conditions.md)
    * [Delegators](docs/handlers/delegators.md)

## License

[MIT License](LICENSE)
