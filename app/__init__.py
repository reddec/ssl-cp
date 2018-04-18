from flask_sqlalchemy import SQLAlchemy
from flask import Flask
from app import config

app = Flask(__name__)
app.config.from_object(config)
db = SQLAlchemy(app)

from . import models
from .routes import *

db.create_all()
