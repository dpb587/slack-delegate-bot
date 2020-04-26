<template>
  <div>
    <If v-if="value.if" v-model="value.if" />
    <Literal v-if="value.literal" v-model="value.literal" />
    <Literalmap v-if="value.literalmap" v-model="value.literalmap" />
    <Pairist v-if="value.pairist" v-model="value.pairist" />
    <Union v-if="value.union" v-model="value.union" />
  </div>
</template>

<script>
import If from './If'
import Literal from './Literal'
import Literalmap from './Literalmap'
import Pairist from './Pairist'
import Union from './Union'

export default {
  components: {
    If,
    Literal,
    Literalmap,
    Pairist,
    Union
  },
  props: {
    value: {
      type: Object,
      required: true
    }
  }
}
</script>
