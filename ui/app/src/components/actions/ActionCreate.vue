<template>
  <div style="display: inline">
    <v-btn text color="primary" @click="showDialog">
      <slot>
        <v-icon>mdi-plus</v-icon>
        new certificate
      </slot>
    </v-btn>
    <v-dialog
        v-model="dialog"
        fullscreen
        hide-overlay
        transition="dialog-bottom-transition"
    >
      <v-card>
        <v-toolbar dark color="primary">
          <v-btn icon dark @click="dialog = false">
            <v-icon>mdi-close</v-icon>
          </v-btn>
          <v-toolbar-title>Issue new certificate</v-toolbar-title>
          <v-spacer></v-spacer>
        </v-toolbar>
        <v-card-text>
          <EditNewCertificate :only-ca="onlyCa" v-model="subject" :issuer="issuer"/>
        </v-card-text>

        <v-card-actions>
          <v-btn color="success" text :loading="creating" @click="create">
            <v-icon>mdi-plus</v-icon>
            create
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import {Emit, Prop} from "vue-property-decorator";
import {Certificate, Subject} from "@/api";
import {API} from "@/store";
import EditNewCertificate from "@/components/EditNewCertificate.vue";

@Component({
  components: {EditNewCertificate}
})
export default class ActionCreate extends Vue {
  @Prop() issuer?: Certificate;
  @Prop({type: Boolean, default: false}) onlyCa!: boolean;

  creating: boolean = false;
  dialog: boolean = false;

  subject: Subject = {}

  showDialog() {
    let subj: Subject = {
      days: 365
    }
    if (this.issuer) {
      subj.issuer = this.issuer.id
    }
    this.subject = subj;
    this.dialog = true;
  }

  @Emit()
  async create(): Promise<Certificate | undefined> {
    this.creating = true
    try {
      let res = await API.createCertificate(this.subject)
      return res.data
    } catch (e) {
      await this.$store.dispatch('error', e)
    } finally {
      this.creating = false;
      this.dialog = false;
    }
  }
}
</script>

<style scoped>

</style>