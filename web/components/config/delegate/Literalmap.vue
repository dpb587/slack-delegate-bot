<template>
  <v-card>
    <v-card-actions>
      <v-list-item>
        <v-list-item-avatar><v-icon>mdi-shuffle</v-icon></v-list-item-avatar>
        <v-list-item-content>
          <v-list-item-title>User Map</v-list-item-title>
          <v-list-item-subtitle>{{ Object.entries(config.users).length }} user{{ Object.entries(config.users).length != 1 ? 's' : '' }}</v-list-item-subtitle>
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
      <v-card-text class="pt-0" v-show="show">
        <div class="mb-4">
          <v-subheader>From</v-subheader>
          <slot name="config-from"></slot>
        </div>

        <div>
          <v-subheader>Users</v-subheader>
          <v-simple-table dense>
            <template v-slot:default>
              <thead>
                <tr>
                  <th class="text-left">Name</th>
                  <th class="text-left">ID</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="(mapping, mappingIdx) in userlist" :key="mapping.name">
                  <td>
                    <v-text-field
                      outlined
                      dense
                      v-model="userlist[mappingIdx].name"
                      required
                    ></v-text-field>
                  </td>
                  <td>
                    <v-text-field
                      outlined
                      dense
                      v-model="userlist[mappingIdx].id"
                      required
                    ></v-text-field>
                  </td>
                </tr>
              </tbody>
            </template>
          </v-simple-table>
        </div>
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
          users: {
            Dan: 'asdf',
            Fores: 'asdfsasdf',
            Gariam: 'ljklk'
          }
        }
      }
    }
  },
  data () {
    return {
      show: this.expanded,
      config: this.value
    }
  },
  computed: {
    userlist: {
      get () {
        return Object.entries(this.config.users).map((v) => {
          return {
            name: v[0],
            id: v[1]
          }
        })
      },
      set (v) {
        this.config.users = v
      }
    }
  }
}
</script>
