<template>
  <v-card>
    <v-card-actions>
      <v-list-item>
        <v-list-item-avatar><img src="https://pair.ist/favicon.ico" /></v-list-item-avatar>
        <v-list-item-content>
          <v-list-item-title>Pair.ist Lookup</v-list-item-title>
          <v-list-item-subtitle>{{ config.team }} {{ config.role || config.track }}</v-list-item-subtitle>
        </v-list-item-content>
      </v-list-item>

      <v-spacer></v-spacer>

      <v-btn
        icon
        @click="show = !show"
      >
        <v-icon>{{ show ? 'mdi-chevron-up' : 'mdi-chevron-down' }}</v-icon>
      </v-btn>
    </v-card-actions>

    <v-expand-transition>
      <v-card-text v-show="show">
        <v-form
            ref="form"
            v-model="valid"
            lazy-validation
          >
            <v-text-field
              outlined
              dense
              v-model="config.team"
              label="Team Name"
              required
            ></v-text-field>

            <v-text-field
              outlined
              dense
              v-model="config.password"
              label="Password"
              type="password"
            ></v-text-field>

            <v-select
              outlined
              dense
              :items="pairist.rolestracks"
              v-model="roletrack"
              label="Role / Track"
            ></v-select>
        </v-form>
      </v-card-text>
    </v-expand-transition>
  </v-card>
</template>

<script>
export default {
  props: {
    expanded: {
      type: Boolean,
      required: false,
      default: true
    },
    value: {
      type: Object,
      required: false,
      default () {
        return {
          team: 'testing',
          password: 'asdfasdf',
          role: 'Interrupt'
        }
      }
    }
  },
  data () {
    return {
      pairist: {
        rolestracks: []
      },
      show: this.expanded,
      config: this.value
    }
  },
  computed: {
    roletrack: {
      get () {
        if (this.config.role) {
          return `role:${this.config.role}`
        } else if (this.config.track) {
          return `track:${this.config.track}`
        }

        return null
      },
      set (v) {
        const value = v.split(':')

        this.config.role = null
        this.config.track = null

        if (value[0] === 'role') {
          this.config.role = value[1].join(':')
        } else if (value[0] === 'track') {
          this.config.track = value[1].join(':')
        }
      }
    }
  }
}
</script>
