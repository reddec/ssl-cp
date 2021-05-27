<template>
  <v-app>

    <v-navigation-drawer v-model="menu" app>
      <v-list dense nav>

        <v-list-item link to="/">
          <v-list-item-icon>
            <v-icon>mdi-home</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>Dashboard</v-list-item-title>
          </v-list-item-content>
          <v-list-item-action>
            <v-list-item-action-text>
              {{ status.total }}
            </v-list-item-action-text>
          </v-list-item-action>
        </v-list-item>

        <v-subheader>CERTIFICATES</v-subheader>

        <v-list-item link :to="{name:'soon-expire'}">
          <v-list-item-icon>
            <v-icon>mdi-clock</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>Soon expire</v-list-item-title>
          </v-list-item-content>
          <v-list-item-action>
            <v-list-item-action-text>
              {{ status.soon_expire }}
            </v-list-item-action-text>
          </v-list-item-action>
        </v-list-item>

        <v-list-item link :to="{name:'expired'}">
          <v-list-item-icon>
            <v-icon>mdi-alarm</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>Expired</v-list-item-title>
          </v-list-item-content>
          <v-list-item-action>
            <v-list-item-action-text>
              {{ status.expired }}
            </v-list-item-action-text>
          </v-list-item-action>
        </v-list-item>


      </v-list>
      <v-list v-if="selected.certificate.id && !selected.certificate.ca" dense nav>
        <v-subheader class="text-truncate">{{ (selected.certificate.name || '').toUpperCase() }}</v-subheader>
        <v-list-item link :to="{name:'certificate', params:{id:selected.certificate.id}}">
          <v-list-item-icon>
            <v-icon>mdi-eye</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>Summary</v-list-item-title>
          </v-list-item-content>
        </v-list-item>

        <v-list-item link :to="{name:'cookbook-nginx-server', params:{id:selected.certificate.id}}">
          <v-list-item-icon>
            <v-icon>mdi-lock</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>Nginx TLS</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </v-list>
    </v-navigation-drawer>

    <v-app-bar flat dense app>
      <v-app-bar-nav-icon @click="menu=!menu"></v-app-bar-nav-icon>
    </v-app-bar>

    <v-main>
      <router-view/>
      <div style="position: absolute; right: 1em; top: 1em; max-width: 40em">
        <transition-group tag="div" name="fade">
          <v-alert
              border="right"
              colored-border
              type="error"
              elevation="2"
              v-for="error in errors"
              :key="error.id"
          >
            {{ error | errMessage }}
          </v-alert>
        </transition-group>
      </div>
    </v-main>
  </v-app>
</template>
<style scoped>
.fade-enter-active, .fade-leave-active {
  transition: opacity .5s;
}

.fade-enter, .fade-leave-to /* .fade-leave-active below version 2.1.8 */
{
  opacity: 0;
}
</style>
<script lang="ts">

import {mapState} from "vuex";
import {CertStore, LocalError} from "@/store";
import {AxiosError} from "axios";
import {Certificate, Status} from "@/api";
import Vue from "vue";
import Component from "vue-class-component";
import {State} from "vuex-class";

@Component({
  filters: {
    errMessage(err: LocalError) {
      let opt = err.error as AxiosError;
      if (!opt.response) {
        return err.error.message;
      }
      if (!opt.response.data) {
        return opt.response.statusText;
      }
      if (!opt.response.data.error) {
        return opt.response.data;
      }
      return opt.response.data.error
    }
  }
})
export default class App extends Vue {

  @State('errors') errors?: LocalError[];

  @State('certificate') selected!: CertStore;

  menu: boolean = !this.mini;

  get status(): Status {
    return this.$store.state.status.status
  }

  get items() {
    return [
      {icon: 'mdi-home', title: 'Dashboard', link: '/'},
      {icon: 'mdi-clock', title: 'Soon expire', link: {name: 'soon-expire'}},
      {icon: 'mdi-alarm', title: 'Expired', link: {name: 'expired'}},
    ]
  }

  get mini() {
    switch (this.$vuetify.breakpoint.name) {
      case 'xs':
        return true
      case 'sm':
        return true
      case 'md':
        return false
      case 'lg':
        return false
      case 'xl':
        return false
    }
  }

  get links() {
    return this.items.map((l) => l.link)
  }
}

</script>
