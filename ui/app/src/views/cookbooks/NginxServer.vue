<template>
  <v-container fluid>

    <v-card :loading="selected.loading">
      <v-card-text>
        <h2>Nginx server configuration</h2>
      </v-card-text>
      <v-card-text>
        <v-switch label="Authorize client certificate"
                  v-model="validateClient"
                  persistent-hint
                  :hint="'Authorize client certificate according to full chain of CA: ' +issuers"
        />
      </v-card-text>
      <v-card-text>
        <h3>Preview config</h3>
        <pre style="overflow-x: auto">{{ nginxConfig }}</pre>
      </v-card-text>
    </v-card>

  </v-container>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import {Getter, State} from "vuex-class";
import CertificatesList from "@/components/CertificatesList.vue";
import {CertStore} from "@/store";
import {Certificate} from "@/api";

@Component({
  components: {CertificatesList}
})
export default class NginxServer extends Vue {
  @State("certificate") selected!: CertStore
  @Getter('certificate/chain') chain!: Certificate[];

  validateClient: boolean = true;

  independent: boolean = false;

  get issuers() {
    return this.chain.map((x) => x.name).join(", ")
  }

  get nginxConfig() {
    const domains = this.selected.certificate.domains?.join(' ') || 'default_server';

    const authorizeFeatures = `
  // authorize client certificate
  ssl_client_certificate  /etc/ssl/${this.selected.certificate.name}/ca.crt;
  ssl_crl                 /etc/ssl/${this.selected.certificate.name}/revoked.crl;
  ssl_verify_client       on;

  // map common name as user name
  map $ssl_client_s_dn $user_name {
    default "";
    ~,CN=(?<CN>[^,]+) $CN;
  }

  // use organization unit as user id - it's equal to unique certificate ID
  map $ssl_client_s_dn $user_id {
    default "";
    ~,OU=(?<OU>[^,]+) $OU;
  }

    `

    const authorizedProxyFeatures = `
    proxy_set_header X-User-ID   $user_id;
    proxy_set_header X-User-Name $user_name;
    proxy_set_header X-Client-DN $ssl_client_s_dn;
    `

    const text = `server {
  listen 443 ssl;
  listen [::]:443 ssl;
  server_name ${domains};

  ssl_certificate         /etc/ssl/${this.selected.certificate.name}/server.crt;
  ssl_certificate_key     /etc/ssl/${this.selected.certificate.name}/server.key;
  ${this.validateClient ? authorizeFeatures : ''}
  location / {
    proxy_set_header X-Real-IP   $remote_addr;
    proxy_set_header Host        $host;
   ${this.validateClient ? authorizedProxyFeatures : ''}
    proxy_pass http://localhost:5000/;
  }
}`

    return text
  }
}
</script>

<style scoped>

</style>