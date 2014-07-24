# -*- coding:utf-8 -*-
"""
Main views
"""

from flask import render_template

from webrss.main import app
from . import category
from . import feed
from . import entry
from webrss.models import Category


@app.route('/')
def index():
    """
    Index view
    """
    categories = Category.select().where(Category.deleted_at == None)
    categories = list(categories)

    return render_template('index.html', categories=categories)


@app.route('/api/search')
def search():
    """
    Return search result
    """
    pass
