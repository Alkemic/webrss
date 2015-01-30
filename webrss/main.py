# -*- coding:utf-8 -*-
"""webrss.main"""
from flask import Flask
from flask import g
from flask_peewee.rest import RestAPI

from webrss.models import DATABASE


app = Flask(__name__)

rest_api = RestAPI(app)


@app.before_request
def peewee_database_connect():
    g.db = DATABASE
    g.db.connect()


@app.after_request
def peewee_database_close(response):
    g.db.close()
    return response
