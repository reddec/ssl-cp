import {Module} from 'vuex'
import type {Certificate, Renewal} from "@/api/api";
import {API, CertStore, RootStore} from "@/store/index";


const store: Module<CertStore, RootStore> = {
    namespaced: true,
    state: {
        loading: false,
        renewing: false,
        certificate: {},
        issued: [],
        chain: []
    },
    getters: {
        ca(state) {
            return !state.loading && state.certificate && state.certificate.ca
        },
        chain(state) {
            return (state.chain || []).filter((x) => state.certificate.ca || x.id !== state.certificate.id)
        }
    },
    mutations: {
        certificate(state, certificate: Certificate) {
            state.certificate = certificate;
        },
        issued(state, list: Certificate[]) {
            state.issued = list
        },
        complete(state) {
            state.loading = false;
        },
        loading(state) {
            state.loading = true;
            state.certificate = {};
            state.issued = [];
            state.chain = [];
        },
        created(state, child: Certificate) {
            if (child.issuer === state.certificate?.id) {
                state.issued = [...state.issued, child]
            }
        },
        chain(state, certificates: Certificate[]) {
            state.chain = certificates;
        },
        renewing(state, value: boolean) {
            state.renewing = value;
        },
        renewed(state, certificate: Certificate) {
            state.renewing = false;
            state.certificate = certificate;
        },
        unselect(state) {
            state.certificate = {}
        }
    },
    actions: {
        async load({commit, state, dispatch}, id: number) {
            if (id === state.certificate.id) {
                return;
            }

            async function walk(id: number): Promise<Certificate[]> {
                const res = await API.getCertificate(id);
                const cert = res.data;
                if (!cert.issuer) {
                    return [cert];
                }
                let up = await walk(cert.issuer);
                return [...up, cert];
            }

            commit('loading')
            try {
                await Promise.all([
                    API.getCertificate(id).then((res) => commit('certificate', res.data)),
                    API.listCertificates(id).then((res) => commit('issued', res.data)),
                    walk(id).then((res) => commit('chain', res)),
                ])

            } catch (e) {
                dispatch('error', e, {root: true})
            } finally {
                commit('complete')
            }
        },
        async renew({commit, dispatch, state}, renewal: Renewal) {
            commit('renewing', true)
            try {
                let newCert = await API.renewCertificate(state.certificate.id || -1, renewal)
                commit('renewed', newCert.data)
            } catch (e) {
                dispatch('error', e, {root: true})
            } finally {
                commit('renewing', false)
            }
        }
    },
    modules: {}
}

export default store;
