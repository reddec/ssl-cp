import Vue from 'vue'
import Vuex from 'vuex'
import type {Certificate, Subject} from "@/api/api";
import {API, LocalError, RootStore} from "@/store";
import CertModule from './cert';
import StatusModule from './status';

Vue.use(Vuex)


const ALERT_TIMEOUT = 3000

let errorID = 0;


export default new Vuex.Store<RootStore>({
    modules: {
        certificate: CertModule,
        status: StatusModule,
    },
    state: {
        loading: false,
        creating: false,
        errors: [],
        certificates: [],
    },
    mutations: {
        loaded(state, certificates: Certificate[]) {
            state.loading = false
            state.certificates = certificates;
        },
        loading(state) {
            state.loading = true;
            state.certificates = [] as Certificate[];
        },
        addError(state, error: LocalError) {
            console.error(error)
            state.errors = [...state.errors, error]
        },
        removeError(state, id: number) {
            state.errors = state.errors.filter((e) => e.id !== id);
        },
        creating(state) {
            state.creating = true;
        },
        created(state, certificate?: Certificate) {
            if (certificate && !state.certificates.find((c) => c.id === certificate.id)) {
                state.certificates = [...state.certificates, certificate]
            }
            state.creating = false;
        }
    },
    actions: {
        async load({commit, dispatch}, parent?: number) {
            commit('loading')
            try {
                let list;
                if (parent) {
                    list = (await API.listCertificates(parent)).data
                } else {
                    list = (await API.listRootCertificates()).data;
                }
                commit('loaded', list)
            } catch (e) {
                commit('loaded', [])
                dispatch('error', e)
            }
        },

        async expired({commit, dispatch}) {
            commit('loading')
            try {
                const list = (await API.listExpiredCertificates()).data
                commit('loaded', list)
            } catch (e) {
                commit('loaded', [])
                dispatch('error', e)
            }
        },

        async soonExpire({commit, dispatch}) {
            commit('loading')
            try {
                const list = (await API.listSoonExpireCertificates()).data
                commit('loaded', list)
            } catch (e) {
                commit('loaded', [])
                dispatch('error', e)
            }
        },

        async preload({state, dispatch}) {
            if (state.certificates && state.certificates.length > 0) {
                return
            }
            dispatch('load')
        },
        error({commit}, error: Error) {
            errorID++;
            const id = errorID;
            commit('addError', {id: id, error: error})
            setTimeout(() => {
                commit('removeError', id)
            }, ALERT_TIMEOUT)
        },
        async create({commit, dispatch, state}, subject: Subject) {
            commit('creating')
            try {
                let res = await API.createCertificate(subject);
                const cert = res.data
                commit('created', cert);
                commit('certificate/created', cert)
            } catch (e) {
                dispatch('error', e)
                commit('created');
            }

        }
    }
})
