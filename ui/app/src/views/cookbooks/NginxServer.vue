<template>
  <v-container fluid>

    <v-card :loading="selected.loading">
      <v-card-text>
        <h2>Nginx server configuration</h2>
        <p>Generator for Nginx server configuration.</p>
      </v-card-text>
      <v-card-text>
        <v-switch label="Authorize client certificate"
                  v-model="config.validateClient"
                  persistent-hint
                  :hint="'Authorize client certificate according to full chain of CA: ' +issuers"
        />
        <v-switch label="Make portable bundle"
                  v-model="config.portable"
                  persistent-hint
                  hint="Create portable bundle that can be used in any location with Nginx"
        />
        <v-text-field
            v-model.number="config.port"
            label="Port"
            hint="binding port incoming connections. Use 8443 to run without root"
            type="number"
        />
      </v-card-text>
      <v-card-text>
        <v-expansion-panels flat>
          <v-expansion-panel>
            <v-expansion-panel-header><h4>Preview config</h4></v-expansion-panel-header>
            <v-expansion-panel-content>
              <pre style="overflow-x: auto">{{ nginxConfig }}</pre>
            </v-expansion-panel-content>
          </v-expansion-panel>
        </v-expansion-panels>
      </v-card-text>
      <v-card-text v-if="config.portable">
        Download archive, extract it, and run server by <br/><kbd>nginx -c $(pwd)/nginx.conf</kbd>
      </v-card-text>
      <v-card-actions>
        <v-btn color="success" text :loading="generating" @click="generate">
          <v-icon>mdi-download</v-icon>
          download
        </v-btn>
      </v-card-actions>
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
import {downloadCertAssets} from "@/lib/utils";
import JSZip from "jszip";
import {saveAs} from "file-saver";

@Component({
  components: {CertificatesList}
})
export default class NginxServer extends Vue {
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

  get sslDir() {
    return this.config.portable ? 'ssl' : `/etc/ssl/${this.selected.certificate.name}`
  }

  get confFile() {
    return this.config.portable ? 'server.conf' : `/etc/nginx/sites-enabled/${this.certificate.id}-${(this.certificate.domains || [this.certificate.name])[0]}.conf`
  }

  async generate() {
    this.generating = true;
    try {
      const assets = await downloadCertAssets(this.selected.chain);


      let zip = new JSZip()
      zip.file(this.sslDir + '/cert.pem', assets.cert)
      zip.file(this.sslDir + '/revoked.pem', assets.revoked.join("\n"))
      if (assets.key) {
        zip.file(this.sslDir + '/key.pem', assets.key)
      }

      let fullchain = assets.cert;

      if (assets.caCerts.length > 0) {
        const ca = assets.caCerts.join('\n')
        fullchain += '\n' + ca
        zip.file(this.sslDir + '/ca.pem', ca)
      }
      if (assets.key) {
        fullchain += '\n' + assets.key
      }

      zip.file(this.sslDir + '/fullchain.pem', fullchain)
      zip.file(this.confFile, this.nginxConfig)
      if (this.config.portable) {
        zip.file('nginx.conf', this.nginxMain)
      }
      const content = await zip.generateAsync({type: 'blob'})
      saveAs(content, (this.certificate.name || 'archive') + '.zip')
    } catch (e) {
      await this.$store.dispatch('error', e)
    } finally {
      this.generating = false;
    }
  }

  get nginxConfig() {
    const domains = this.selected.certificate.domains?.join(' ') || 'default_server';

    const authorizeFeatures = `
  # authorize client certificate
  ssl_client_certificate  ${this.sslDir}/ca.pem;
  ssl_crl                 ${this.sslDir}/revoked.pem;
  ssl_verify_client       on;

    `

    const authorizedProxyFeatures = `
    proxy_set_header X-Client-DN $ssl_client_s_dn;
    `

    return `server {
  listen ${this.config.port} ssl;
  listen [::]:${this.config.port} ssl;
  server_name ${domains};

  ssl_certificate         ${this.sslDir}/cert.pem;
  ssl_certificate_key     ${this.sslDir}/key.pem;
  ${this.config.validateClient ? authorizeFeatures : ''}
  location / {
    proxy_set_header X-Real-IP   $remote_addr;
    proxy_set_header Host        $host;
   ${this.config.validateClient ? authorizedProxyFeatures : ''}
    proxy_pass http://localhost:5000/;
  }
}`

  }

  get nginxMain() {
    return `
worker_processes auto;
daemon off;
pid /tmp/nginx-${this.certificate.id}.pid;
error_log stderr info;
events {
        worker_connections 768;
}

http {
        sendfile on;
        tcp_nopush on;
        tcp_nodelay on;
        keepalive_timeout 65;
        types_hash_max_size 2048;
        default_type application/octet-stream;
        ssl_protocols TLSv1 TLSv1.1 TLSv1.2 TLSv1.3; # Dropping SSLv3, ref: POODLE
        ssl_prefer_server_ciphers on;

        access_log /dev/stdout;
        error_log stderr;
        gzip on;
        include server.conf;
}
`
  }
}
</script>

<style scoped>

</style>