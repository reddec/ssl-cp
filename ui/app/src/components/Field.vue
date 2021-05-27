<template>
  <v-row class="wrap" :class="klass">
    <v-col cols="10" sm="4" md="3" lg="2" xl="2" v-text="label"
           class="font-weight-bold"></v-col>
    <v-col cols="2" sm="1" md="1" lg="1" xl="1">
      <slot name="icon">
        <v-icon>{{defaultIcon}}</v-icon>
      </slot>
    </v-col>
    <v-col cols="12" sm="7" md="8" lg="9" xl="9">
      <slot/>
    </v-col>

  </v-row>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import {Prop} from "vue-property-decorator";

@Component
export default class Field extends Vue {
  @Prop({default: ''}) label!: string;

  @Prop({default: false, type: Boolean}) warning!: boolean;
  @Prop({default: false, type: Boolean}) failed!: boolean;

  get klass() {
    if (this.failed) {
      return 'error'
    }
    if (this.warning) {
      return 'warning'
    }
    return '';
  }

  get defaultIcon() {
    if (this.failed) {
      return 'mdi-alert'
    }
    if (this.warning) {
      return 'mdi-alert-circle-outline'
    }
    return ''
  }
}
</script>

<style scoped>

</style>