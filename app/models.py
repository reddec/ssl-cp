from app import db
from datetime import datetime
from sqlalchemy.orm import deferred


class Project(db.Model):
    id = db.Column(db.Integer, primary_key=True)
    title = db.Column(db.Text, nullable=False, unique=True)
    description = db.Column(db.Text, nullable=False, default=lambda: '')
    ca_private = deferred(db.Column(db.Binary))
    ca_cert = deferred(db.Column(db.Binary))


class Certificate(db.Model):
    id = db.Column(db.Integer, primary_key=True)
    project_id = db.Column(db.Integer, db.ForeignKey('project.id'), nullable=False)
    private_key = deferred(db.Column(db.Binary, nullable=False))
    public_cert = deferred(db.Column(db.Binary, nullable=False))
    created_at = db.Column(db.DateTime, default=datetime.utcnow, nullable=False)
    expire_at = db.Column(db.DateTime, nullable=False)
    common_name = db.Column(db.Text, nullable=False, unique=True)
    revoked_at = db.Column(db.DateTime)

    project = db.relationship('Project', backref=db.backref('certificates', lazy=True))

    @property
    def is_active(self):
        n = datetime.now()
        return n < self.expire_at and self.revoked_at is None
