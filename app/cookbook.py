from app import app, db, models, api
from flask import make_response
import tarfile, io
import os


@app.route('/cookbook/stunnel/<cert_id>/server')
def stunnel_server(cert_id: int):
    cert = models.Certificate.query.get(int(cert_id))  # type: models.Certificate
    project = cert.project  # type: models.Project
    certs = models.Certificate.revoked(cert.project_id)
    ca_file = project.ca_cert  # type: bytes
    key_file = cert.private_key
    cert_file = cert.public_cert
    crl_file = api.create_revoke_list(ca_file, key_file, [(cert.id, cert.revoked_at) for cert in certs])

    with open(os.path.join(os.path.dirname(__file__), "assets", "cookbook", "stunnel", "config.conf"), 'rb') as f:
        config = f.read()

    mem_arch = io.BytesIO()
    dirname = cert.common_name
    with tarfile.open(fileobj=mem_arch, mode='w:gz') as tar:
        add_to_archive(tar, config, dirname + '/stunnel.conf')
        add_to_archive(tar, ca_file, dirname + '/ca.cert')
        add_to_archive(tar, key_file, dirname + '/node.key')
        add_to_archive(tar, cert_file, dirname + '/node.cert')
        add_to_archive(tar, crl_file, dirname + '/crl.pem')
    mem_arch.seek(0)
    resp = make_response(mem_arch.read())
    resp.headers['Content-Type'] = 'application/gzip'
    resp.headers['Content-Disposition'] = 'attachment; filename="' + str(cert.id) + "-" + cert.common_name + ".tar.gz"
    return resp


def add_to_archive(tar: tarfile.TarFile, buf: bytes, name: str):
    info = tarfile.TarInfo(name)
    info.size = len(buf)
    info.type = tarfile.REGTYPE
    tar.addfile(info, io.BytesIO(buf))
