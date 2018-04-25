from flask import render_template, make_response
from app.utils.archive import dict_to_archive
from app.routes.cookbook import render
from app import app, models


@app.route('/certificate/<cert_id>/nginx')
def nginx_cookbook(cert_id: int):
    cert = models.Certificate.query.get(int(cert_id))
    return render_template('cookbook/nginx.html', project=cert.project, cert=cert)


@app.route('/certificate/<cert_id>/nginx/basic')
def nginx_basic_config(cert_id: int):
    cert = models.Certificate.query.get(int(cert_id))  # type: models.Certificate
    template = render('cookbook', 'nginx', 'nginx.conf', cert=cert)
    readme = render('cookbook', 'nginx', 'README.md.j2', cert=cert, project=cert.project)
    archive = dict_to_archive({
        cert.common_name + '/README.md': readme,
        cert.common_name + '/nginx.conf': template,
        cert.common_name + '/' + str(cert_id) + ".crt": cert.public_cert,
        cert.common_name + '/' + str(cert_id) + ".key": cert.private_key,
    })
    resp = make_response(archive)
    resp.headers['Content-Type'] = 'application/gzip'
    resp.headers['Content-Disposition'] = 'attachment; filename="nginx-' + cert.common_name + "-basic.tar.gz"
    return resp


@app.route('/certificate/<cert_id>/nginx/with-auth')
def nginx_auth_config(cert_id: int):
    cert = models.Certificate.query.get(int(cert_id))  # type: models.Certificate
    project = cert.project
    template = render('cookbook', 'nginx', 'nginx-auth.conf', cert=cert, project=project)
    readme = render('cookbook', 'nginx', 'README.md.j2', cert=cert, project=project)
    archive = dict_to_archive({
        cert.common_name + '/README.md': readme,
        cert.common_name + '/nginx.conf': template,
        cert.common_name + '/ca-' + str(project.id) + ".crt": project.ca_cert,
        cert.common_name + '/' + str(cert_id) + ".crt": cert.public_cert,
        cert.common_name + '/' + str(cert_id) + ".key": cert.private_key,
    })
    resp = make_response(archive)
    resp.headers['Content-Type'] = 'application/gzip'
    resp.headers['Content-Disposition'] = 'attachment; filename="nginx-' + cert.common_name + "-auth.tar.gz"
    return resp
