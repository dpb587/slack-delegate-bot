<template>
  <v-card>
    <v-card-actions>
      <v-list-item>
        <v-list-item-avatar><v-icon>mdi-shuffle</v-icon></v-list-item-avatar>
        <v-list-item-content>
          <v-list-item-title>Conditional</v-list-item-title>
          <v-list-item-subtitle>&nbsp;</v-list-item-subtitle>
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
          <v-subheader>If</v-subheader>
          TODO
        </div>

        <div>
          <v-subheader>Then</v-subheader>
          <AbstractDelegate v-model="config.then" />
          {{ config.then }}
        </div>
      </v-card-text>
    </v-expand-transition>
  </v-card>
</template>

<script>
import AbstractDelegate from './Abstract'

export default {
  components: {
    AbstractDelegate
  },
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
        return {}
      }
    }
  },
  data () {
    return {
      show: this.expanded,
      config: this.value
    }
  }
}
</script>
