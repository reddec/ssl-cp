import os

SQLALCHEMY_DATABASE_URI = os.getenv('DATABASE_URI', 'sqlite:////tmp/test.db')
FLASK_DEBUG = 1
