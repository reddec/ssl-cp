from typing import List

from app import app, models
from datetime import datetime
from flask import request
from sqlalchemy import or_


@app.route('/api/project/<id>/certificates', methods=['GET'])
def api_list_certificates(id: int):
    certs = models.Certificate.query.filter(models.Certificate.project_id == int(id))  # type: List[models.Certificate]
    return "\n".join(str(cert.id) for cert in certs)


@app.route('/api/project/<id>/certificates', methods=['POST'])
def api_create_certificat(id: int):
    dbcert = models.create_child(int(id), request.form['cn'], int(request.form['days']))
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
