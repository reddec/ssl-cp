from typing import List, Tuple
from datetime import datetime, timedelta

from OpenSSL import crypto, SSL

digest = 'sha256'


def create_ca(common_name: str, days: int) -> (bytes, bytes):
    pair = crypto.PKey()
    pair.generate_key(crypto.TYPE_RSA, 2048)

    cert = crypto.X509()
    cert.get_subject().CN = common_name
    cert.gmtime_adj_notBefore(0)
    cert.gmtime_adj_notAfter(days * 24 * 60 * 60)
    cert.set_issuer(cert.get_subject())
    cert.set_pubkey(pair)
    cert.sign(pair, digest)
    return crypto.dump_certificate(crypto.FILETYPE_PEM, cert), \
           crypto.dump_privatekey(crypto.FILETYPE_PEM, pair)


def create_cert(cn: str, days: int, ca_private: bytes, ca_cert: bytes, sn: int) -> (bytes, bytes):
    pkey = crypto.PKey()
    pkey.generate_key(crypto.TYPE_RSA, 2048)

    req = crypto.X509Req()
    req.get_subject().CN = cn
    req.set_pubkey(pkey)
    req.sign(pkey, digest)

    issuerKey = crypto.load_privatekey(crypto.FILETYPE_PEM, ca_private)
    issuerCert = crypto.load_certificate(crypto.FILETYPE_PEM, ca_cert)

    cert = crypto.X509()
    cert.gmtime_adj_notBefore(0)
    cert.gmtime_adj_notAfter(days * 24 * 60 * 60)
    cert.set_issuer(issuerCert.get_subject())
    cert.set_subject(req.get_subject())
    cert.set_pubkey(req.get_pubkey())
    cert.set_serial_number(sn)
    cert.sign(issuerKey, digest)

    return crypto.dump_certificate(crypto.FILETYPE_PEM, cert), \
           crypto.dump_privatekey(crypto.FILETYPE_PEM, pkey)


def create_revoke_list(ca_cert: bytes, ca_key: bytes, serials: List[Tuple[int, datetime]]) -> bytes:
    crl = crypto.CRL()

    crl.set_lastUpdate(datetime.utcnow().strftime('%Y%m%d%H%M%SZ').encode())
    crl.set_nextUpdate((datetime.utcnow() + timedelta(days=365)).strftime('%Y%m%d%H%M%SZ').encode())
    for serial, revoked_at in serials:
        revoked = crypto.Revoked()
        revoked.set_serial(hex(serial)[2:].encode())
        revoked.set_reason(b'keyCompromise')
        revoked.set_rev_date(revoked_at.strftime('%Y%m%d%H%M%SZ').encode())
        crl.add_revoked(revoked)
    key = crypto.load_privatekey(crypto.FILETYPE_PEM, ca_key)
    cert = crypto.load_certificate(crypto.FILETYPE_PEM, ca_cert)
    crl.sign(cert, key, digest.encode())
    return crypto.dump_crl(crypto.FILETYPE_PEM, crl)
