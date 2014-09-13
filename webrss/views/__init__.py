# -*- coding:utf-8 -*-
"""
Main views
"""
import os

from flask import render_template
from flask import request
from flask import send_from_directory

from webrss.main import app
from . import category
from . import feed
from . import entry
from webrss.models import Category, Entry


@app.route('/')
def index():
    """
    Index view
    """
    categories = Category.select().where(Category.deleted_at.__eq__(None))
    categories = list(categories)

    return render_template('index.html', categories=categories)


@app.route('/api/search', methods=['POST'])
def search():
    """
    Return search result
    """
    phrase = '%%%s%%' % request.form['phrase']

    entries = Entry.select()\
        .where((Entry.title ** phrase) | (Entry.summary ** phrase))\
        .where(Entry.deleted_at.__eq__(None))

    return render_template('search.html', entries=list(entries))


@app.route('/favicon.ico')
def favicon():
    return send_from_directory(
        os.path.join(app.root_path, 'static'),
        'favicon.ico',
        mimetype='image/png'
    )