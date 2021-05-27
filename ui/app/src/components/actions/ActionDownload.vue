<template>
  <v-btn text color="primary" @click="download" :loading="generating">
   <v-icon>mdi-download</v-icon> download
  </v-btn>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import {Emit, Prop} from "vue-property-decorator";
import {API} from "../../store";
import {Certificate} from "../../api";
import JSZip from "jszip";
import {saveAs} from 'file-saver';

@Component
export default class ActionDownload extends Vue {

  @Prop() certificate!: Certificate;
  @Prop() chain?: Certificate[];

  generating: boolean = false;

  @Emit()
  async download() {
    this.generating = true;
    try {
      const [
        publicCert,
        caCerts,
      ] = await Promise.all([
        API.getPublicCert(this.certificate.id!).then((x) => x.data),
        Promise.all((this.chain || [])
            .filter((c) => c.id !== this.certificate.id)
            .map((cert) => API.getPublicCert(cert.id!).then((x) => x.data))),
      ])
      let privateKey;
      if (!this.certificate.ca) {
        privateKey = (await API.getPrivateKey(this.certificate.id!)).data;
      }

      let zip = new JSZip()
      zip.file('cert.pem', publicCert)
      if (privateKey) {
        zip.file('key.pem', privateKey)
      }

      let fullchain = publicCert;

      if (caCerts.length > 0) {
        const ca = caCerts.join('\n')
        fullchain += '\n' + ca
        zip.file('ca.pem', ca)
      }
      if (privateKey) {
        fullchain += '\n' + privateKey
      }

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