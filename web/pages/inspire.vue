<template>
  <v-layout>
    <Abstract v-model="config.delegatebot.delegate" />
  </v-layout>
</template>

<script>
import Abstract from '~/components/config/delegate/Abstract'

export default {
  components: {
    Abstract
  },
  data () {
    return {
      config: {
        delegatebot: {
          watch: [
            {
              target: {
                channel: 'C02HWMDUQ'
              }
            }
          ],
          delegate: {
            if: {
              when: [
                {
                  or: [
                    {
                      and: [
                        {
                          hours: {
                            tz: 'America/New_York',
                            start: '09:00',
                            end: '18:00'
                          }
                        },
                        {
                          day: {
                            tz: 'America/New_York',
                            days: [
                              'Mon',
                              'Tue',
                              'Wed',
                              'Thu',
                              'Fri'
                            ]
                          }
                        },
                        {
                          not: {
                            tz: 'America/New_York',
                            dates: [
                              '2019-01-01',
                              '2019-01-21'
                            ]
                          }
                        }
                      ]
                    }
                  ]
                }
              ],
              then: {
                literalmap: {
                  from: {
                    union: [
                      {
                        literal: {
                          text: 'Interrupts are logged'
                        }
                      },
                      {
                        pairist: {
                          team: 'cfbuildpacks',
                          role: 'Interrupt'
                        }
                      },
                      {
                        pairist: {
                          team: 'cfbuildpacks-releng',
                          role: 'Interrupt'
                        }
                      }
                    ]
                  },
                  users: {
                    Dan: 'UCL6FS0HW',
                    Fores: 'U1PMZB0TT',
                    Garima: 'U2LAQS9RQ'
                  }
                }
              }
            }
          },
          options: {
            empty_message: "Sorry, we couldn't magically find the buildpacks interrupt."
          }
        }
      }
    }
  }
}
</script>
