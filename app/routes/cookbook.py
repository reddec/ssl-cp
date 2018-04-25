from app import app, db, models, certs_tools
from app.utils.archive import add_to_archive
from flask import make_response, request
from jinja2 import Environment, FileSystemLoader
import tarfile, io
import os


def render(*asset_path, **kwargs):
    env = Environment(loader=FileSystemLoader(os.path.join(os.path.dirname(__file__), '..', "assets")))
    templ = env.get_template(os.path.join(*asset_path))
    return templ.render(**kwargs).encode()


@app.route('/cookbook/stunnel/<cert_id>/server')
def stunnel_server(cert_id: int):
    cert = models.Certificate.query.get(int(cert_id))  # type: models.Certificate
    project = cert.project  # type: models.Project
    certs = models.Certificate.revoked(cert.project_id)
    ca_cert = project.ca_cert  # type: bytes
    ca_key = project.ca_private
    key_file = cert.private_key
    cert_file = cert.public_cert
    crl_file = certs_tools.create_revoke_list(ca_cert, ca_key, [(cert.id, cert.revoked_at) for cert in certs])
    config = render('cookbook', 'stunnel', 'server.conf', cert=cert,
                    accept=request.args.get('accept', 0),
                    connect=request.args.get('connect', 0))
    mem_arch = io.BytesIO()
    dirname = cert.common_name + "-server"
    with tarfile.open(fileobj=mem_arch, mode='w:gz') as tar:
        add_to_archive(tar, config, dirname + '/stunnel.conf')
        add_to_archive(tar, ca_cert, dirname + '/ca.cert')
        add_to_archive(tar, key_file, dirname + '/node.key')
        add_to_archive(tar, cert_file, dirname + '/node.cert')
        add_to_archive(tar, crl_file, dirname + '/crl.pem')
        add_to_archive(tar, render('cookbook', 'stunnel', 'server.service', cert=cert),
                       dirname + '/stunnel-' + cert.common_name + '-server.service')
        add_to_archive(tar, render('cookbook', 'stunnel', 'server-install.sh', cert=cert),
                       dirname + '/install.sh')
    mem_arch.seek(0)
    resp = make_response(mem_arch.read())
    resp.headers['Content-Type'] = 'application/gzip'
    resp.headers['Content-Disposition'] = 'attachment; filename="' + str(
        cert.id) + "-" + cert.common_name + "-server.tar.gz"
    return resp


@app.route('/cookbook/stunnel/<cert_id>/client')
def stunnel_client(cert_id: int):
    cert = models.Certificate.query.get(int(cert_id))  # type: models.Certificate
    project = cert.project  # type: models.Project
    certs = models.Certificate.revoked(cert.project_id)
    ca_cert = project.ca_cert  # type: bytes
    ca_key = project.ca_private
    key_file = cert.private_key
    cert_file = cert.public_cert
    crl_file = certs_tools.create_revoke_list(ca_cert, ca_key, [(cert.id, cert.revoked_at) for cert in certs])
    config = render('cookbook', 'stunnel', 'client.conf', cert=cert,
                    accept=request.args.get('accept', 0),
                    connect=request.args.get('connect', 0))
    mem_arch = io.BytesIO()
    dirname = cert.common_name + "-client"
    with tarfile.open(fileobj=mem_arch, mode='w:gz') as tar:
        add_to_archive(tar, config, dirname + '/stunnel.conf')
        add_to_archive(tar, ca_cert, dirname + '/ca.cert')
        add_to_archive(tar, key_file, dirname + '/node.key')
        add_to_archive(tar, cert_file, dirname + '/node.cert')
        add_to_archive(tar, crl_file, dirname + '/crl.pem')
        add_to_archive(tar, render('cookbook', 'stunnel', 'client.service', cert=cert),
                       dirname + '/stunnel-' + cert.common_name + '-client.service')
        add_to_archive(tar, render('cookbook', 'stunnel', 'client-install.sh', cert=cert),
                       dirname + '/install.sh')
    mem_arch.seek(0)
    resp = make_response(mem_arch.read())
    resp.headers['Content-Type'] = 'application/gzip'
    resp.headers['Content-Disposition'] = 'attachment; filename="' + str(
        cert.id) + "-" + cert.common_name + "-client.tar.gz"
    return resp
