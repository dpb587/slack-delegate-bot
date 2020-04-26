<template>
  <v-card>
    <v-card-actions>
      <v-list-item>
        <v-list-item-avatar><v-icon>mdi-text</v-icon></v-list-item-avatar>
        <v-list-item-content>
          <v-list-item-title>Literal Text</v-list-item-title>
          <v-list-item-subtitle>{{ config.text.substr(0, 48) }}{{ config.text.length > 48 ? '...' : '' }}</v-list-item-subtitle>
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
            <v-textarea
              outlined
              dense
              v-model="config.text"
              name="input-7-4"
              label="Text"
            ></v-textarea>
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
          text: 'Interrupts are logged here - https://drive.google.com/open?id=1vk8dSnRKGirkHb0vCbaZaTtTNDHVsFQkJCmOrvtbSIs'
        }
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
