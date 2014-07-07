# -*- coding:utf-8 -*-
from flask import render_template, request, g
from peewee import fn
from webrss.decorators import jsonify
from webrss.functions import categories_dict
from webrss.main import app
from webrss.models import Category, Feed


@app.route('/api/category/list', endpoint='category.index')
@jsonify
def index():
    return {'categories': categories_dict()}


@app.route('/api/category/create', endpoint='category.create', methods=['POST'])
@jsonify
def create():
    order_max = Category.select(fn.Max(Category.order).alias('max_order'))[0].max_order
    order_max = 0 if order_max is None else order_max
    try:
        Category.create(title=request.form['category-name'], order=order_max + 1)

        return {'status': 'ok'}
    except Exception as e:
        return {'status': 'fail'}


@app.route('/api/category/update/<int:pk>', endpoint='category.update', methods=['POST'])
@jsonify
def update(pk):
    category = Category.get(Category.id == pk)

    return {}


@app.route('/api/category/delete/<int:pk>', endpoint='category.delete', methods=['POST'])
@jsonify
def delete(pk):
    entries = Category.select().where(Category.deleted_at == None).dicts()

    return {'categories': categories_dict()}
