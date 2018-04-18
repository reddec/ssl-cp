from app import app, models, db, certs_tools
from flask import request, url_for, render_template, redirect, abort, make_response


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
    cert, priv = certs_tools.create_ca(request.form['cn'], int(request.form['days']))
    project.ca_private = priv
    project.ca_cert = cert
    db.session.commit()
    return redirect(url_for('get_project', id=project.id))


@app.route('/project/<id>/generate-signed', methods=['POST'])
def create_child_certificate(id: int):
    crt = models.create_child(int(id), request.form['cn'], int(request.form['days']))
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
    content = certs_tools.create_revoke_list(project.ca_cert, project.ca_private,
                                             [(cert.id, cert.revoked_at) for cert in certs])
    resp = make_response(content)
    resp.headers['Content-Type'] = 'text/plain'
    resp.headers['Content-Disposition'] = 'attachment; filename="revoked-' + str(id) + ".crl"
    return resp


@app.route('/project/<project_id>/certificates')
def project_certificates(project_id: int):
    proj = models.Project.query.get(int(project_id))
    return render_template('certificates.html', project=proj)
