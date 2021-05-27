<template>
  <div>
    <v-chip class="ma-1" outlined small label color="warning" v-if="status.soonExpire">soon expire</v-chip>
    <v-chip class="ma-1" small label color="warning" v-if="status.expired">expired</v-chip>
    <v-chip class="ma-1" small label color="primary" outlined v-if="certificate.ca">CA</v-chip>
  </div>
</template>

<script lang="ts">

import Component from "vue-class-component";
import Vue from "vue";
import {Certificate} from "@/api";
import {Prop} from "vue-property-decorator";
import Status from "@/lib/status";

@Component
export default class CertificateBadges extends Vue {
  @Prop() certificate!: Certificate

  get status() {
    return new Status(this.certificate)
  }
}

</script>

<style scoped>

</style>