<template>
  <div>
    <v-row>
      <CertificateBadges :certificate="certificate"/>
    </v-row>
    <v-row>
      <v-col cols="12" md="6">
        <Field label="Created">
          {{ certificate.created_at | longDate }} ({{ certificate.created_at | ago }})
        </Field>
        <Field label="Updated">
          {{ certificate.updated_at | longDate }} ({{ certificate.updated_at | ago }})
        </Field>
        <Field label="Duration">
          {{ status.duration }} days
        </Field>
        <Field label="Expire" :warning="status.soonExpire" :failed="status.expired">
          {{ certificate.expire_at | longDate }} ({{ certificate.expire_at | ago }})
        </Field>
      </v-col>
      <v-col cols="12" md="6">
        <Field label="Role">
          {{ certificate.ca ? 'Central Authority' : 'Client Auth' }}
        </Field>
        <Field label="Serial">
          {{ certificate.serial }}
        </Field>
        <Field label="Issuer">
          <router-link :to="{name:'certificate', params:{id:issuer.id}}" v-if="issuer">
            {{ issuer.name }}
          </router-link>
          <span v-else>N/A</span>
        </Field>
        <Field label="Domains">
          <div v-if="!certificate.ca">
            {{ (certificate.domains || []).join(', ') }}
          </div>
          <span v-else>N/A</span>
        </Field>

      </v-col>
    </v-row>
    <v-row v-if="certificate.units">
      <v-col>
        <Field label="Organization units">
          {{ (certificate.units || []).join(', ') }}
        </Field>
      </v-col>
    </v-row>
  </div>
</template>

<script lang="ts">
import type {Certificate} from "@/api";
import Component from 'vue-class-component'
import Vue, {PropType} from "vue";
import dayjs from "dayjs";
import CertificateBadges from "@/components/CertificateBadges.vue";
import Field from "@/components/Field.vue";
import {Prop} from "vue-property-decorator";
import Status from "@/lib/status";

@Component({
  components: {Field, CertificateBadges},
})
export default class CertificateInfo extends Vue {
  @Prop() certificate!: Certificate;
  @Prop() chain?: Certificate[];
  @Prop({type: Boolean, default: false}) loading!: boolean;

  get status() {
    return new Status(this.certificate)
  }

  get issuer() {
    return this.chain && this.chain.length > 1 ? this.chain[this.chain.length - 2] : undefined
  }
}
</script>

<style scoped>

</style>