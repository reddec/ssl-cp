<template>
  <div style="display: inline-block">
    <v-btn text color="error" :loading="revoking" @click="showDialog">revoke</v-btn>
    <v-dialog
        v-model="dialog"
        max-width="30em"
    >
      <v-card>
        <v-card-text class="pt-5">
          <h3>Are you sure?</h3>
        </v-card-text>
        <v-card-text class="pt-5">
          It will also revoke all nested certificates recursively
        </v-card-text>
        <v-card-actions>
          <v-btn color="primary" text :loading="revoking" @click="dialog = false">
            cancel
          </v-btn>
          <v-spacer/>
          <v-btn color="error" text :loading="revoking" @click="revoke">
            <v-icon>mdi-alert</v-icon>
            revoke
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
import {Certificate} from "@/api";
import {API} from "@/store";
import dayjs from "dayjs";

@Component
export default class ActionRevoke extends Vue {
  @Prop() readonly certificate!: Certificate;

  revoking: boolean = false;
  dialog: boolean = false;

  showDialog() {
    this.dialog = true;
  }

  @Emit()
  async revoke(): Promise<any> {
    this.revoking = true
    try {
      await API.revokeCertificate(this.certificate.id || -1)
    } catch (e) {
      await this.$store.dispatch('error', e)
    } finally {
      this.revoking = false;
      this.dialog = false;
    }
  }
}
</script>

<style scoped>

</style>