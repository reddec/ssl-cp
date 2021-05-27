<template>
  <v-data-table
      :headers="headers"
      item-key="id"
      :sort-by="['ca', 'name']"
      :search="search"
      :items="items"
      :items-per-page="full ? -1 : 20"
      :item-class="rowColor"
      @click:row="openCertificate"
  >
    <template v-slot:top>
      <v-text-field
          v-model="search"
          label="Search"
          class="mx-4"
      >
      </v-text-field>
    </template>
    <template v-slot:item.ca="{item}">
      {{ item.ca ? 'Central Authority' : 'Client' }}
    </template>
    <template v-slot:item.domains="{item}">
      <span v-if="!item.ca">{{ (item.domains || []).join(', ') }}</span>
      <span v-else>N/A</span>
    </template>
    <template v-slot:item.created_at="{item}">
      {{ item.created_at | ago }}
    </template>
    <template v-slot:item.updated_at="{item}">
      {{ item.updated_at | ago }}
    </template>
    <template v-slot:item.expire_at="{item}">
      {{ item.expire_at | ago }}
    </template>
  </v-data-table>
</template>

<script lang="ts">
import type {Certificate} from "@/api";
import Vue from "vue";
import Component from "vue-class-component";
import {Prop} from "vue-property-decorator";
import Status from "@/lib/status";

@Component
export default class CertificatesList extends Vue {
  @Prop({default: []}) readonly certificates!: Certificate[];

  @Prop({default: false, type: Boolean}) readonly full!: boolean;

  @Prop(Boolean) readonly onlyCa!: boolean;

  search: string = '';

  get headers() {
    let headers = [{text: "Name", value: "name"},]
    if (!this.onlyCa) {
      headers.push({text: "Role", value: "ca"})
    }
    return [...headers,
      {text: "Domains", value: "domains"},
      {text: "Created", value: "created_at"},
      {text: "Updated", value: "updated_at"},
      {text: "Expire", value: "expire_at"},
    ]
  }

  get items() {
    if (this.onlyCa) {
      return (this.certificates || []).filter((c) => c.ca)
    }
    return this.certificates || [];
  }

  openCertificate(item?: Certificate) {
    if (!item) return
    this.$router.push({name: 'certificate', params: {id: item.id + ''}})
  }

  rowColor(item: Certificate) {
    let status = new Status(item);
    if (status.expired) {
      return 'error'
    }
    if (status.soonExpire) {
      return 'warning'
    }
    return ''
  }
}

</script>

<style scoped>

</style>