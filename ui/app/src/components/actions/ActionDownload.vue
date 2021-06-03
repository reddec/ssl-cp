<template>
  <v-btn text color="primary" @click="download" :loading="generating">
    <slot>
      <v-icon>mdi-download</v-icon>
      download
    </slot>
  </v-btn>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import {Emit, Prop} from "vue-property-decorator";
import {Certificate} from "../../api";
import JSZip from "jszip";
import {saveAs} from 'file-saver';
import {downloadCertAssets} from "@/lib/utils";

@Component
export default class ActionDownload extends Vue {

  @Prop() certificate!: Certificate;
  @Prop() chain?: Certificate[];

  generating: boolean = false;

  @Emit()
  async download() {
    this.generating = true;
    try {
      const assets = await downloadCertAssets(this.chain || [this.certificate]);


      let zip = new JSZip()
      zip.file('cert.pem', assets.cert)
      zip.file('revoked.pem', assets.revoked.join("\n"))
      if (assets.key) {
        zip.file('key.pem', assets.key)
      }

      let fullchain = assets.cert;

      if (assets.caCerts.length > 0) {
        const ca = assets.caCerts.join('\n')
        fullchain += '\n' + ca
        zip.file('ca.pem', ca)
      }
      if (assets.key) {
        fullchain += '\n' + assets.key
      }

      fullchain += '\n' + assets.revoked.join("\n")

      zip.file('fullchain.pem', fullchain)
      const content = await zip.generateAsync({type: 'blob'})
      saveAs(content, (this.certificate.name || 'archive') + '.zip')
    } catch (e) {
      await this.$store.dispatch('error', e)
    } finally {
      this.generating = false;
    }
  }

}
</script>

<style scoped>

</style>