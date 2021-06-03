import Vue from 'vue'
import VueRouter from 'vue-router'
import Home from '../views/Home.vue'
import Certificate from '../views/certificate/Index.vue'
import CertificateLayout from "@/views/certificate/layout.vue";
import Expired from "@/views/Expired.vue";
import SoonExpire from "@/views/SoonExpire.vue";
import NginxServer from "@/views/cookbooks/NginxServer.vue";
import Clients from "@/views/cookbooks/Clients.vue";
Vue.use(VueRouter)

const routes = [
    {
        path: '/',
        name: 'home',
        component: Home
    },
    {
        path: '/certificate/:id',
        component: CertificateLayout,
        props: true,
        children: [
            {
                path: '/',
                name: 'certificate',
                component: Certificate
            },
            {
                path: 'cookbook/nginx-server',
                name: 'cookbook-nginx-server',
                component: NginxServer
            },
            {
                path: 'cookbook/clients',
                name: 'cookbook-clients',
                component: Clients
            }
        ]
    },
    {
        path: '/certificates/expired',
        name: 'expired',
        component: Expired,
    },
    {
        path: '/certificates/soon-expire',
        name: 'soon-expire',
        component: SoonExpire,
    }
]

const router = new VueRouter({
    routes
})

export default router
