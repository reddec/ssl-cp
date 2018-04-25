from typing import List

from app import db
from datetime import datetime, timedelta
from sqlalchemy.orm import deferred


class Project(db.Model):
    id = db.Column(db.Integer, primary_key=True)
    title = db.Column(db.Text, nullable=False, unique=True)
    description = db.Column(db.Text, nullable=False, default=lambda: '')
    ca_private = deferred(db.Column(db.Binary))
    ca_cert = deferred(db.Column(db.Binary))

    certificates = db.relationship('Certificate', back_populates='project', lazy=True)


class Certificate(db.Model):
    id = db.Column(db.Integer, primary_key=True)
    project_id = db.Column(db.Integer, db.ForeignKey('project.id'), nullable=False)
    private_key = deferred(db.Column(db.Binary, nullable=False))
    public_cert = deferred(db.Column(db.Binary, nullable=False))
    created_at = db.Column(db.DateTime, default=datetime.utcnow, nullable=False)
    expire_at = db.Column(db.DateTime, nullable=False)
    common_name = db.Column(db.Text, nullable=False, unique=True)
    revoked_at = db.Column(db.DateTime)

    project = db.relationship('Project', back_populates='certificates')

    @property
    def is_active(self):
        n = datetime.now()
        return n < self.expire_at and self.revoked_at is None

    @staticmethod
    def revoked(project_id: int) -> List['Certificate']:
        return Certificate.query \
            .filter(Certificate.project_id == project_id) \
            .filter(Certificate.revoked_at.isnot(None))


def create_child(project_id: int, cn: str, days: int) -> Certificate:
    from app import certs_tools
    init_at = datetime.now()
    expire_at = init_at + timedelta(days=days)
    project = Project.query.get(project_id)
    cert, key = b'', b''
    dbcert = Certificate(project_id=project_id,
                         common_name=cn,
                         expire_at=expire_at,
                         private_key=key,
                         public_cert=cert)

    db.session.add(dbcert)  # get SN
    db.session.flush()
    cert, key = certs_tools.create_cert(cn, days, project.ca_private, project.ca_cert, int(dbcert.id))
    dbcert.public_cert = cert
    dbcert.private_key = key
    db.session.commit()
    return dbcert
