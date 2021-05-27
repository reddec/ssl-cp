<template>
  <div>
    <p class="text-h5 text-center text-truncate">{{ label }}</p>
    <div class="text-h1  text-center">
      <slot/>
    </div>
    <p class="text-caption text-center text-truncate" style="min-height: 2em">{{ hint }}</p>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import {Prop} from "vue-property-decorator";

@Component
export default class MetricNumber extends Vue {
  @Prop({default: ''}) label!: string;
  @Prop({default: ''}) hint!: string;
  @Prop({default: ''}) color!: string;

  get textColor(){
    return `text--${this.color}`
  }
}
</script>

<style scoped>

</style>