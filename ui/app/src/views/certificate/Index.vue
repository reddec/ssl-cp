<template>
  <v-container>
    <v-row class="flex-wrap">
      <v-col cols="12">
        <v-card>
          <v-card-title>
            <h2>{{ selected.certificate.name }}</h2>
          </v-card-title>
          <v-card-text>
            <CertificateInfo :certificate="selected.certificate" :chain="selected.chain"/>
          </v-card-text>
          <v-card-actions>
            <ActionRenew @renew="certUpdated" :certificate="selected.certificate"/>
            <ActionDownload :certificate="selected.certificate" :chain="selected.chain"/>
            <v-spacer/>
            <ActionRevoke @revoke="certRemoved(selected.certificate)" :certificate="selected.certificate"/>
          </v-card-actions>
        </v-card>
      </v-col>
    </v-row>

    <v-row v-if="selected.certificate.ca">
      <v-col cols="12">
        <v-card>
          <v-card-title>Issued certificates</v-card-title>
          <v-card-text>
            <CertificatesList :certificates="selected.issued"/>
          </v-card-text>
          <v-card-actions>
            <action-create :issuer="selected.certificate" @create="certCreated"/>
          </v-card-actions>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts">
import CertificateInfo from "@/components/CertificateInfo.vue";
import CertificatesList from "@/components/CertificatesList.vue";
import EditNewCertificate from "@/components/EditNewCertificate.vue";
import {Certificate as ApiCert, Subject} from "@/api";
import Vue from "vue";
import {Component, Prop, Watch} from "vue-property-decorator";
import {State} from "vuex-class";
import {CertStore} from "@/store";
import ActionRenew from "@/components/actions/ActionRenew.vue";
import ActionRevoke from "@/components/actions/ActionRevoke.vue";
import ActionCreate from "@/components/actions/ActionCreate.vue";
import ActionDownload from "@/components/actions/ActionDownload.vue";
import Field from "@/components/Field.vue";


@Component({
  components: {
    Field,
    ActionDownload,
    ActionCreate, ActionRevoke, ActionRenew, EditNewCertificate, CertificatesList, CertificateInfo
  }
})
export default class Certificate extends Vue {

  certCreated(certificate?: ApiCert) {
    if (certificate) {
      this.$store.commit('created', certificate)
      this.$store.commit('certificate/created', certificate)
    }
  }

  certUpdated(certificate?: ApiCert) {
    if (certificate) {
      this.$store.commit('certificate/certificate', certificate)
      this.$store.dispatch('status/load')
    }
  }

  certRemoved(certificate?: ApiCert) {
    this.$store.dispatch('status/load')
    if (!certificate || !certificate.issuer) {
      this.$router.push({name: 'home'})
      return
    }
    this.$router.push({name: 'certificate', params: {id: certificate.issuer.toString()}})
  }

  @State("certificate") selected!: CertStore

  @State('creating') creating!: boolean;

  @State(state => state.certificate.renewing) renewing!: boolean;

}
</script>

<style scoped>

</style>