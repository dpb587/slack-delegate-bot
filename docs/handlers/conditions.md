# Conditions

Conditions are used to require context-specific details for something can occur.


## `and`

A list of conditionals, all of which must be satisfied.

```yaml
and:
- day: { days: [ Mon, Tue, Wed, Thu, Fri ] }
- hours: { start: 09:00, end: 18:00 }
```


## `date`

Satisfied when the current date matches one of the listed dates (`YYYY-MM-DD` format).

```yaml
date:
  tz: America/Los_Angeles
  dates: [ 2019-01-01, 2019-01-21, 2019-02-18, 2019-05-27 ]
```


## `day`

Satisfied when the current day matches one of the listed days (`Mon` format).

```yaml
day:
  tz: Europe/Berlin
  days: [ Mon, Tue, Wed, Thu, Fri ]
```


## `hours`

Satisfied when the current time is within a start and end time (`HH:MM` 24h format).

```yaml
hours:
  tz: America/Toronto
  start: 09:00
  end: 17:00
```


## `not`

Inverts another conditional's result.

```yaml
not:
  day: { days: [ Mon, Tue, Wed, Thu, Fri ] }
```


## `target`

Satisfied when targeting a specific channel.

```yaml
target:
  channel: C02HPPYQ2
```


## `or`

A list of conditionals, one of which must be satisfied.

```yaml
or:
- hours: { start: 09:00, end: 12:30 }
- hours: { start: 13:30, end: 17:00 }
```
