from app import app, models, db
from flask import render_template, make_response, request, redirect, url_for
from datetime import datetime


@app.route('/certificate/<cert_id>/stunnel')
def stunnel_cookbook(cert_id: int):
    cert = models.Certificate.query.get(int(cert_id))
    return render_template('stunnel.html', project=cert.project, cert=cert)


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
