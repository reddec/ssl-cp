import {API} from "@/store";
import {Certificate} from "@/api";
import cert from "@/store/cert";

export interface Asset {
    chain: Certificate[]
    cert: string
    key?: string
    caCerts: string[]
    revoked: string[]
}

export async function downloadCertAssets(chain: Certificate[]): Promise<Asset> {
    const last = chain.length - 1
    let certificate = chain[last];

    let fetchPK: Promise<{ data: string }>;
    if (!certificate.ca) {
        fetchPK = API.getPrivateKey(certificate.id!);
    } else {
        fetchPK = new Promise((resolve, reject) => resolve({data: ''}))
    }

    const [
        publicCert,
        caCerts,
        privateKey,
        crls,
    ] = await Promise.all([
        API.getPublicCert(certificate.id!).then((x) => x.data),
        Promise.all(
            (chain || [])
                .slice(0, last)
                .map((cert) => API.getPublicCert(cert.id!).then((x) => x.data))
        ),
        fetchPK.then((r) => r.data),
        Promise.all(
            (chain || [])
                .filter((x) => x.ca)
                .map((cert) => API.getRevokedCertificatesList(cert.id!).then((x) => x.data))
        ),
    ])

    return {
        cert: publicCert,
        key: privateKey,
        caCerts: caCerts,
        chain: chain,
        revoked: crls
    }
}