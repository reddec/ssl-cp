from app import app, models
from flask import render_template


@app.route('/project/<project_id>/certificates')
def all_certificates(project_id: int):
    proj = models.Project.query.get(int(project_id))
    return render_template('certificates.html', project=proj)


@app.route('/certificate/<cert_id>/stunnel')
def stunnel_cookbook(cert_id: int):
    cert = models.Certificate.query.get(int(cert_id))
    return render_template('stunnel.html', project=cert.project, cert=cert)
