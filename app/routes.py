from typing import List

from app import app, db, models
from app import api, cookbook
from app.certificates import *

from sqlalchemy import or_
from flask import render_template, request, redirect, url_for, make_response, abort
from datetime import datetime, timedelta

db.create_all()


@app.route('/')
@app.route('/projects')
def projects():
    projects = models.Project.query.all()
    return render_template('projects.html', projects=projects)


@app.route('/projects', methods=['POST'])
def add_project():
    name = request.form['name']
    proj = models.Project(title=name)
    db.session.add(proj)
    db.session.commit()
    return redirect(url_for('get_project', id=proj.id))


@app.route('/project/<id>', methods=['POST'])
def update_project(id: int):
    id = int(id)
    project = models.Project.query.get(id)
    project.description = request.form['description']
    db.session.commit()
    return redirect(url_for('get_project', id=project.id))


@app.route('/project/<id>', methods=['GET'])
def get_project(id: int):
    id = int(id)
    project = models.Project.query.get(id)
    return render_template('project.html', project=project)


@app.route('/project/<id>/generate-root', methods=['POST'])
def gen_root_ca(id: int):
    id = int(id)
    project = models.Project.query.get(id)
    cert, priv = api.create_ca(request.form['cn'], int(request.form['days']))
    project.ca_private = priv
    project.ca_cert = cert
    db.session.commit()
    return redirect(url_for('get_project', id=project.id))


def create_child(project_id: int, cn: str, days: int) -> models.Certificate:
    init_at = datetime.now()
    expire_at = init_at + timedelta(days=days)
    project = models.Project.query.get(project_id)
    if project.ca_private is None or project.ca_cert is None:
        return abort(412)
    cn = cn.strip()
    if cn == '':
        return abort(412)
    cert, key = b'', b''
    dbcert = models.Certificate(project_id=project_id,
                                common_name=cn,
                                expire_at=expire_at,
                                private_key=key,
                                public_cert=cert)

    db.session.add(dbcert)  # get SN
    db.session.flush()
    cert, key = api.create_cert(cn, days, project.ca_private, project.ca_cert, int(dbcert.id))
    dbcert.public_cert = cert
    dbcert.private_key = key
    db.session.commit()
    return dbcert


@app.route('/project/<id>/generate-signed', methods=['POST'])
def create_child_certificate(id: int):
    crt = create_child(int(id), request.form['cn'], int(request.form['days']))
    return redirect(url_for('certificate', id=crt.id))


@app.route('/project/<id>/root/cert')
def download_project_ca_cert(id: int):
    id = int(id)
    project = models.Project.query.get(id)
    if project.ca_cert is None:
        return abort(404)
    resp = make_response(project.ca_cert)
    resp.headers['Content-Type'] = 'text/plain'
    resp.headers['Content-Disposition'] = 'attachment; filename="ca-' + str(project.id) + ".crt"
    return resp


@app.route('/project/<id>/revoked')
def download_revoked_crl(id: int):
    id = int(id)
    certs = models.Certificate.revoked(id)
    project = models.Project.query.get(id)  # type: models.Project
    content = api.create_revoke_list(project.ca_cert, project.ca_private,
                                     [(cert.id, cert.revoked_at) for cert in certs])
    resp = make_response(content)
    resp.headers['Content-Type'] = 'text/plain'
    resp.headers['Content-Disposition'] = 'attachment; filename="revoked-' + str(id) + ".crl"
    return resp


@app.route('/certificate/<id>')
def certificate(id: int):
    id = int(id)
    cert = models.Certificate.query.get(id)  # type: models.Certificate
    return render_template('certificate.html', cert=cert, project=cert.project)


@app.route('/certificate/<id>/cert')
def download_certificate_cert(id: int):
    id = int(id)
    cert = models.Certificate.query.get(id)  # type: models.Certificate
    resp = make_response(cert.public_cert)
    resp.headers['Content-Type'] = 'text/plain'
    resp.headers['Content-Disposition'] = 'attachment; filename="' + str(id) + ".crt"
    return resp


@app.route('/certificate/<id>/key')
def download_certificate_key(id: int):
    id = int(id)
    cert = models.Certificate.query.get(id)  # type: models.Certificate
    resp = make_response(cert.private_key)
    resp.headers['Content-Type'] = 'text/plain'
    resp.headers['Content-Disposition'] = 'attachment; filename="' + str(id) + ".key"
    return resp


@app.route('/revoked', methods=['POST'])
def revoke_certificate():
    cert_id = int(request.form['cert'])
    crt = models.Certificate.query.get(cert_id)  # type: models.Certificate
    crt.revoked_at = datetime.now()
    db.session.commit()
    return redirect(url_for('certificate', id=cert_id))


@app.route('/api/project/<id>/certificates', methods=['GET'])
def api_list_certificates(id: int):
    certs = models.Certificate.query.filter(models.Certificate.project_id == int(id))  # type: List[models.Certificate]
    return "\n".join(str(cert.id) for cert in certs)


@app.route('/api/project/<id>/certificates', methods=['POST'])
def api_create_certificat(id: int):
    dbcert = create_child(int(id), request.form['cn'], int(request.form['days']))
    return str(dbcert.id)


@app.route('/api/project/<id>/certificates/active')
def api_list_active_certificates(id: int):
    certs = models.Certificate.query \
        .filter(models.Certificate.project_id == int(id)) \
        .filter(models.Certificate.expire_at >= datetime.now()) \
        .filter(models.Certificate.revoked_at.is_(None))  # type: List[models.Certificate]
    return "\n".join(str(cert.id) for cert in certs)


@app.route('/api/project/<id>/certificates/inactive')
def api_list_project_inactive_certificates(id: int):
    certs = models.Certificate.query \
        .filter(models.Certificate.project_id == int(id)) \
        .filter(or_(
        models.Certificate.expire_at < datetime.now(),
        models.Certificate.revoked_at.isnot(None)
    ))  # type: List[models.Certificate]
    return "\n".join(str(cert.id) for cert in certs)
