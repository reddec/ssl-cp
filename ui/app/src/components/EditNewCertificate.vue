<template>
  <div class="mt-4">
    <v-row>
      <v-col>
        <v-card v-if="issuer" flat outlined>
          <v-card-title>Issuer</v-card-title>
          <v-card-text>
            Will be signed by {{ issuer.name }}
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
    <v-row>
      <v-col>
        <v-card flat outlined>
          <v-card-title>New certificate</v-card-title>
          <v-card-text>
            <v-text-field label="Common name" hint="/CN in certificate"
                          :value="local.name"
                          :disabled="renew"
                          @input="patch('name', $event)"/>
            <v-text-field
                type="number"
                label="Duration in days" hint="Certificate will be expired after defined days" :value="local.days"
                @input="patch('days', parseInt($event))"/>
          </v-card-text>
          <v-card-text>
            <v-switch
                label="Certificate Authority"
                persistent-hint
                :disabled="renew || onlyCa"
                :input-value="isCA"
                @change="patch('ca', $event)"
                hint="Can issue another certificates but can not be used for client-server authorization"/>
          </v-card-text>
          <v-card-text>
            <list-edit
                title="Organization Units"
                label="Unit name"
                hint="will be added as /OU in the certificate"
                icon="mdi-tag"
                :value="local.units"
                @change="patch('units', $event)"/>
          </v-card-text>
          <v-card-text v-if="!isCA">
            <list-edit
                title="Domains (Subject Alternative Name)"
                label="Domain name"
                hint="will be added as /SAN in certificate subject"
                icon="mdi-domain"
                :value="domains"
                @change="patch('domains', $event)"/>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

  </div>
</template>

<script lang="ts">
import {Certificate, Subject} from "../api";
import {Component, Prop, Vue} from "vue-property-decorator";
import ListEdit from "@/components/ListEdit.vue";

@Component({
  components: {ListEdit}
})
export default class EditNewCertificate extends Vue {
  @Prop()
  readonly value!: Subject

  @Prop()
  readonly issuer?: Certificate


  @Prop({default: false, type: Boolean})
  readonly onlyCa!: boolean;


  @Prop({default: false, type: Boolean})
  readonly renew!: boolean;

  patch(name: string, value: any) {
    this.$emit('input', {...this.local, [name]: value})
  }

  get local() {
    return this.value || {}
  }

  get domains() {
    return this.local.domains || []
  }

  get isCA() {
    return this.onlyCa || this.local.ca
  }

}

</script>

<style scoped>

</style>