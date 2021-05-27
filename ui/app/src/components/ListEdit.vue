<template>
  <v-list>
    <h4 v-text="title"></h4>
    <v-list-item>
      <v-list-item-content>
        <v-form @submit.prevent="add">
          <v-text-field
              :label="label"
              :hint="hint"
              v-model="newValue"
              append-outer-icon="mdi-add"
          />
        </v-form>
      </v-list-item-content>
    </v-list-item>
    <v-list-item
        v-for="value in local"
        :key="value"
    >
      <v-list-item-avatar>
        <v-icon
            class="grey lighten-1"
            dark
            v-text="icon"
        >
        </v-icon>
      </v-list-item-avatar>

      <v-list-item-content>
        <v-list-item-title v-text="value"></v-list-item-title>
      </v-list-item-content>

      <v-list-item-action>
        <v-btn icon @click="remove(value)">
          <v-icon color="grey lighten-1">mdi-delete</v-icon>
        </v-btn>
      </v-list-item-action>
    </v-list-item>
  </v-list>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import {Model, Prop} from "vue-property-decorator";

@Component
export default class DomainsList extends Vue {
  @Prop({default: ''}) readonly title!: string;
  @Prop({default: ''}) readonly label!: string;
  @Prop({default: ''}) readonly hint!: string;
  @Prop({default: 'mdi-info'}) readonly icon!: string;
  @Model('change') readonly value!: string[];

  newValue: string = '';

  add() {
    if (this.local.findIndex((x) => x === this.newValue) == -1) {
      this.$emit('change', [...this.local, this.newValue]);
    }
    this.newValue = '';
  }

  remove(value: string) {
    this.$emit('change', [...this.local.filter((d) => d !== value)]);
  }

  get local() {
    return this.value || [];
  }
}
</script>

<style scoped>

</style>