# -*- coding:utf-8 -*-
from datetime import datetime
from flask import render_template, request

from webrss.decorators import jsonify
from webrss.main import app
from webrss.models import Category, Feed, Entry


@app.route('/api/entry/fetch', endpoint='entry.fetch', methods=['POST'])
def fetch():
    entry = Entry.get(Entry.id == request.form['pk'])
    entry.read_at = datetime.now()
    entry.save()

    return render_template('entry/fetch.html', entry=entry)
