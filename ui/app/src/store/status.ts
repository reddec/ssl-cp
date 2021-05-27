import {Module} from "vuex";
import {API, RootStore, StatusStore} from "@/store/index";
import {Status} from "@/api";

const store: Module<StatusStore, RootStore> = {
    namespaced: true,
    state: {
        loading: false,
        status: {}
    },
    mutations: {
        status(state, status?: Status) {
            state.status = status || {};
            state.loading = false;
        },
        loading(state) {
            state.loading = true;
        }
    },
    actions: {
        async load({commit, dispatch}) {
            commit('loading')
            try {
                const res = await API.getStatus();
                commit('status', res.data)
            } catch (e) {
                commit('status')
                dispatch('error', e, {root: true})
            }
        },
    }
}

export default store;