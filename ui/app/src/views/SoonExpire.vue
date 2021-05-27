<template>
  <v-container>
    <v-row>
      <v-col>
        <v-card :loading="loading">
          <v-card-text>
            <h2>Soon expire certificates</h2>
          </v-card-text>
          <v-card-text>
            <CertificatesList :certificates="certificates"/>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import {State} from "vuex-class";
import {Certificate} from "@/api";
import CertificatesList from "@/components/CertificatesList.vue";
@Component({
  components: {CertificatesList}
})
export default class SoonExpire extends Vue {
  @State('certificates') readonly certificates!: Certificate[];
  @State('loading') readonly loading!: boolean;

  beforeMount() {
    this.$store.dispatch('soonExpire')
  }
}
</script>

<style scoped>

</style>