# Slack Setup

This bot relies on its own [Slack App](https://api.slack.com/apps) with the following settings.

## Settings

### App Home

 * **Always Show My Bot as Online** -- enabled
 * **Home Tab** -- disabled
 * **Messages Tab** -- enabled

### OAuth & Permissions

 * **Bot Token Scopes**
    * **app_mentions:read**
    * **chat:write**
    * **im:history**
    * **mpim:history**
    * **users:read**
 * **Restrict API Token Usage**

### Event Subscriptions

Only enable this after the bot is deployed and running.

 * **Events** -- enabled
 * **Request URL** -- `{botURL}/api/v1/slack/event`
 * **Bot Events**
    * **app_mention**
    * **message.im**
    * **message.mpim**

### User ID Translation

 * **Translate Global IDs** -- disabled
