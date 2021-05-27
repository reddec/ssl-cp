<template>
  <v-container>
    <v-row>
      <v-col>
        <v-card :loading="statusLoading">
          <v-card-text>
            <metric-text label="Total" hint="issued certificates">
              {{ status.total || 0 }}
            </metric-text>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col class="d-none d-xl-block">
        <v-card :loading="statusLoading">
          <v-card-text>
            <metric-text label="CA" hint="Central Authorities">
              {{ status.ca || 0 }}
            </metric-text>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col>
        <v-card :to="{name:'expired'}" :color="status.expired > 0 ? 'error' : ''" :loading="statusLoading">
          <v-card-text>
            <metric-text label="Expired" hint="outdated but not revoked">
              {{ status.expired || 0 }}
            </metric-text>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col>
        <v-card :to="{name:'soon-expire'}" :color="status.soon_expire > 0 ? 'warning' : ''" :loading="statusLoading">
          <v-card-text>
            <metric-text label="Soon expire" hint="expire within 30 days">
              {{ status.soon_expire || 0 }}
            </metric-text>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col>
        <v-card :loading="statusLoading">
          <v-card-text>
            <metric-text label="Live" hint="all not revoked">
              {{ (status.total || 0) - (status.expired || 0) }}
            </metric-text>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col class="d-none d-md-block">
        <v-card :loading="statusLoading">
          <v-card-text>
            <metric-text label="Revoked">
              {{ status.revoked || 0 }}
            </metric-text>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
    <v-row>
      <v-col>
        <v-card :loading="loading">
          <v-card-title class="flex align-content-center">
            <h2>Root certificates</h2>
            <ActionCreate @create="created" only-ca v-if="count > 5">
              <v-icon>mdi-plus</v-icon>
              add
            </ActionCreate>
          </v-card-title>
          <v-card-text>
            <CertificatesList :certificates="certificates" full only-ca/>
          </v-card-text>
          <v-card-actions>
            <ActionCreate @create="created" only-ca/>
          </v-card-actions>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts">
import Vue from 'vue';
import CertificatesList from "@/components/CertificatesList.vue";
import {mapState} from "vuex";
import ActionCreate from "@/components/actions/ActionCreate.vue";
import {Certificate, Status} from "@/api";
import MetricText from "@/components/MetricText.vue"; // @ is an alias to /src

export default Vue.extend({
  name: 'Home',
  components: {
    MetricText,
    ActionCreate,
    CertificatesList,
  },
  beforeMount() {
    this.$store.dispatch('load');
    this.$store.dispatch('status/load')
  },
  computed: {
    ...mapState(['loading', 'certificates']),
    count() {
      return (this.certificates || []).length
    },
    status(): Status {
      return this.$store.state.status.status
    },
    statusLoading(): boolean {
      return this.$store.state.status.loading;
    }
  },
  methods: {
    created(certificate: Certificate) {
      this.$store.commit('created', certificate)
    }
  }
});
</script>
