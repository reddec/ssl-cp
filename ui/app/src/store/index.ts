import {Certificate, Configuration, DefaultApi, Status} from "@/api";
import {AxiosError} from "axios";

export const API = new DefaultApi(new Configuration({
    basePath: process.env.VUE_APP_API
}))

export interface LocalError {
    id: number,
    error: Error | AxiosError
}

export interface CertStore {
    loading: boolean,
    renewing: boolean,
    certificate: Certificate,
    issued: Certificate[],
    chain: Certificate[]
}

export interface RootStore {
    loading: boolean,
    creating: boolean,
    errors: LocalError[],
    certificates: Certificate[]
}

export interface StatusStore {
    loading: boolean,
    status: Status
}