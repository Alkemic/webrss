# -*- coding:utf-8 -*-
from flask import render_template, json

from webrss.decorators import jsonify
from webrss.functions import categories_dict, categories_tuple
from webrss.main import app
import category
import feed
import entry
from webrss.models import Category


@app.route('/')
def index():
    categories = Category.select().where(Category.deleted_at == None)

    return render_template('index.html', categories=[entry for i, entry in enumerate(categories)])


@app.route('/api/search')
def search():
    pass
