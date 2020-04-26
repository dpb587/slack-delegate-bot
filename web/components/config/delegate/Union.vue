<template>
  <v-card>
    <v-card-actions>
      <v-list-item>
        <v-list-item-avatar><v-icon>mdi-filter-none</v-icon></v-list-item-avatar>
        <v-list-item-content>
          <v-list-item-title>Combine</v-list-item-title>
          <v-list-item-subtitle>TODO</v-list-item-subtitle>
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
        <slot name="config-todo"></slot>
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
        return []
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
