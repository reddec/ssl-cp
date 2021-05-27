import Vue from 'vue'
import App from './App.vue'
import router from './router'
import store from './store/root'
import vuetify from './plugins/vuetify'
import 'roboto-fontface/css/roboto/roboto-fontface.css'
import '@mdi/font/css/materialdesignicons.css'
import dayjs, {Dayjs} from "dayjs";
import relativeTime from 'dayjs/plugin/relativeTime';

dayjs.extend(relativeTime)

Vue.config.productionTip = false

Vue.filter('ago', (str: string) => dayjs(str).fromNow())
Vue.filter('shortDate', (value: Dayjs | string) => dayjs(value).format('DD MMM YY'))
Vue.filter('longDate', (value: Dayjs | string) => dayjs(value).format('ddd, DD MMM YYYY HH:mm:ss'))

new Vue({
    router,
    store,
    vuetify,
    render: function (h) {
        return h(App)
    }
}).$mount('#app')
