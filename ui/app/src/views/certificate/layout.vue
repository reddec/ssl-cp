<template>
  <div>
    <v-breadcrumbs :items="breadcrumbs"></v-breadcrumbs>
    <v-progress-linear indeterminate v-if="selected.loading"/>
    <router-view/>
  </div>
</template>

<script lang="ts">
import CertificateInfo from "@/components/CertificateInfo.vue";
import CertificatesList from "@/components/CertificatesList.vue";
import EditNewCertificate from "@/components/EditNewCertificate.vue";
import Vue from "vue";
import {Component, Prop, Watch} from "vue-property-decorator";
import {State} from "vuex-class";
import {CertStore} from "@/store";
import ActionRenew from "@/components/actions/ActionRenew.vue";
import ActionRevoke from "@/components/actions/ActionRevoke.vue";
import ActionCreate from "@/components/actions/ActionCreate.vue";
import ActionDownload from "@/components/actions/ActionDownload.vue";

Component.registerHooks([
  'beforeRouteLeave',
])

@Component({
  components: {
    ActionDownload,
    ActionCreate, ActionRevoke, ActionRenew, EditNewCertificate, CertificatesList, CertificateInfo
  }
})
export default class CertificateLayout extends Vue {
  @Prop() id!: string;

  beforeMount() {
    this.$store.dispatch('certificate/load', parseInt(this.id))
    this.$store.dispatch('preload');
  }

  @Watch('id')
  idChanged() {
    this.$store.dispatch('certificate/load', parseInt(this.id))
  }

  get breadcrumbs() {
    return this.selected.chain.map((cert) => {
      return {
        text: cert.name,
        to: {name: 'certificate', params: {id: cert.id}},
      }
    })
  }

  beforeRouteLeave(to: any, from: any, next: () => any) {
    this.$store.commit('certificate/unselect')
    next()
  }

  @State("certificate") selected!: CertStore

  @State('creating') creating!: boolean;

  @State(state => state.certificate.renewing) renewing!: boolean;

}
</script>

<style scoped>

</style>