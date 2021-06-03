<template>
  <v-container fluid>
    <v-row>
      <v-col>
        <v-card :loading="selected.loading">
          <v-card-text>
            <h2>Quick start</h2>
          </v-card-text>
          <v-card-text>
            <ol>
              <li>Download
                <ActionDownload :certificate="selected.certificate" :chain="chain">archive</ActionDownload>
              </li>
              <li>Unpack</li>
              <li>Use example for your language</li>
            </ol>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
    <v-row>
      <v-col>
        <v-card>
          <v-card-text>
            <h2>Curl</h2>
          </v-card-text>
          <v-card-text>
            <pre>curl --cert fullchain.pem --cacert ca.pem --crlfile revoked.pem https://{{ hostname }}/</pre>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import {Getter, State} from "vuex-class";
import CertificatesList from "@/components/CertificatesList.vue";
import {CertStore} from "@/store";
import {Certificate} from "@/api";
import {downloadCertAssets} from "@/lib/utils";
import JSZip from "jszip";
import {saveAs} from "file-saver";
import ActionDownload from "@/components/actions/ActionDownload.vue";

@Component({
  components: {ActionDownload, CertificatesList}
})
export default class Clients extends Vue {
  @State("certificate") selected!: CertStore
  @Getter('certificate/chain') chain!: Certificate[];
  @State((state) => state.certificate.certificate) certificate!: Certificate;


  config = {
    validateClient: true,
    port: 8443,
    portable: false,
  }

  generating: boolean = false;

  get issuers() {
    return this.chain.map((x) => x.name).join(", ")
  }

  get unc() {
    return this.certificate.domains && this.certificate.domains.length;
  }


  get hostname() {
    if (this.certificate.domains && this.certificate.domains.length > 0) {
      return this.certificate.domains[0];
    }
    if (this.certificate.ips && this.certificate.ips.length > 0) {
      return this.certificate.ips[0];
    }
    return this.certificate.name
  }

}
</script>

<style scoped>

</style>