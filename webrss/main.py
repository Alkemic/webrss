# -*- coding:utf-8 -*-
from flask import Flask
from flask import g

from models import database


app = Flask(__name__)


@app.before_request
def peewee_database_connect():
    g.db = database
    g.db.connect()


@app.after_request
def peewee_database_close(response):
    g.db.close()
    return response
