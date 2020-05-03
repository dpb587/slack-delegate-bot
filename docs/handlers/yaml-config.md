# YAML Configuration

The bot may be configured through one or more YAML configuration files. There must be a top-level `delegatebot` key which contains your configuration and inside it should contain the following keys:

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
