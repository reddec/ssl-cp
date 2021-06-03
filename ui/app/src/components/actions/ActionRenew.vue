<template>
  <div style="display: inline-block">
    <v-btn text color="success" :loading="renewing" @click="showDialog">renew</v-btn>
    <v-dialog
        v-model="dialog"
        :fullscreen="!certificate.ca"
        :hide-overlay="!certificate.ca"
        max-width="40em"
    >
      <v-card>
        <v-toolbar
            dark
            color="primary"
        >
          <v-btn
              icon
              dark
              @click="dialog = false"
          >
            <v-icon>mdi-close</v-icon>
          </v-btn>
          <v-toolbar-title>Renew certificate</v-toolbar-title>
          <v-spacer></v-spacer>
        </v-toolbar>
        <v-card-text class="pt-5">
          <v-text-field
              type="number"
              label="Duration in days" hint="Certificate will be expired after defined days"
              v-model.number="renewal.days"/>

        </v-card-text>
        <v-card-text>
          <list-edit
              title="Organization Units"
              label="Unit name"
              hint="will be added as /OU in the certificate"
              icon="mdi-tag"
              v-model="renewal.units"/>
        </v-card-text>
        <v-card-text v-if="!certificate.ca">
          <list-edit
              title="IP addresses"
              label="IP"
              hint="will be added as /SAN IP in the certificate"
              icon="mdi-ip"
              v-model="renewal.ips"/>
          <list-edit
              title="Domains (Subject Alternative Name)"
              label="Domain name"
              hint="will be added as /SAN DNS in certificate subject"
              icon="mdi-domain"
              v-model="renewal.domains"/>
        </v-card-text>
        <v-card-actions>
          <v-btn color="success" text :loading="renewing" @click="renew">
            <v-icon>mdi-restart</v-icon>
            renew
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
import {Certificate, Renewal} from "@/api";
import EditNewCertificate from "@/components/EditNewCertificate.vue";
import {API} from "@/store";
import dayjs from "dayjs";
import ListEdit from "@/components/ListEdit.vue";

@Component({
  components: {ListEdit, EditNewCertificate}
})
export default class ActionRenew extends Vue {
  @Prop() readonly certificate!: Certificate;

  renewing: boolean = false;
  renewal: Renewal = {
    days: 0,
    domains: [],
    units: [],
    ips: [],
  }
  dialog: boolean = false;

  showDialog() {
    this.renewal.domains = [...(this.certificate.domains || [])];
    this.renewal.units = [...(this.certificate.units || [])];
    this.renewal.ips = [...(this.certificate.ips || [])];
    this.renewal.days = dayjs(this.certificate.expire_at).diff(this.certificate.updated_at || '', 'days');
    this.dialog = true;
  }

  @Emit()
  async renew(): Promise<Certificate | undefined> {
    this.renewing = true
    try {
      const newCert = await API.renewCertificate(this.certificate.id || -1, this.renewal)
      return newCert.data
    } catch (e) {
      await this.$store.dispatch('error', e)
    } finally {
      this.renewing = false;
      this.dialog = false;
    }
  }
}
</script>

<style scoped>

</style>