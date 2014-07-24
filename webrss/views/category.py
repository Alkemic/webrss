# -*- coding:utf-8 -*-
"""
Category related views
"""

from datetime import datetime

from flask import request
from peewee import fn, DoesNotExist, DatabaseError

from webrss.decorators import jsonify
from webrss.functions import categories_dict
from webrss.main import app
from webrss.models import Category


@app.route('/api/category/list', endpoint='category.index')
@jsonify
def index():
    return {'categories': categories_dict()}


@app.route(
    '/api/category/create',
    endpoint='category.create',
    methods=['POST'],
)
@jsonify
def create():
    """
    Create new category
    """
    order_max = Category.select(fn.Max(Category.order).alias('max_order'))[
        0].max_order
    order_max = 0 if order_max is None else order_max
    try:
        Category.create(title=request.form['category-name'],
                        order=order_max + 1)

        return {'status': 'ok'}
    except Exception as e:
        return {'status': 'fail'}


@app.route(
    '/api/category/update',
    endpoint='category.update',
    methods=['POST', 'GET'],
)
@jsonify
def update():
    """
    Update existing category
    """
    is_get = request.method == 'GET'
    pk = request.args.get('pk', None) if is_get else request.form.get('pk',
                                                                      None)

    if pk is None:
        return {'status': 'fail', 'message': 'Wrong parameter'}

    try:
        entry = Category.get(Category.id == pk)
        """ :type : Category """
    except ValueError:
        return {'status': 'fail', 'message': 'Wrong parameter'}
    except DoesNotExist:
        return {'status': 'fail', 'message': 'Entry doesn\'t exists'}

    if is_get:  # if requesting via GET, return feed data
        return {'name': entry.title}

    if not all(name in request.form for name in ('pk', 'category-name')):
        return {'status': 'fail', 'message': 'Not all parameters were sent'}

    try:
        entry.title = request.form['category-name']
        entry.save()
        return {'status': 'ok'}
    except DatabaseError:
        return {'status': 'fail', 'message': 'Exception occurred during save'}


@app.route(
    '/api/category/delete',
    endpoint='category.delete',
    methods=['POST'],
)
@jsonify
def delete():
    """
    Delete category
    """
    try:
        pk = int(request.form['pk'])
        entry = Category.get(Category.id == pk)
    except ValueError:
        return {'status': 'fail', 'message': 'Wrong parameter'}
    except DoesNotExist:
        return {'status': 'fail', 'message': 'Entry doesn\'t exists'}

    entry.deleted_at = datetime.now()

    try:
        entry.save()
    except DatabaseError:
        return {
            'status': 'fail',
            'message': 'Exception occurred during deleting'
        }

    return {'status': 'ok'}


@app.route(
    '/api/category/move_up',
    endpoint='category.move_up',
    methods=['POST'],
)
@jsonify
def move_up():
    """
    Move category up
    """
    try:
        pk = int(request.form['pk'])
        entry = Category.get(Category.id == pk)
        """:type : Category """
    except ValueError:
        return {'status': 'fail', 'message': 'Wrong parameter'}
    except DoesNotExist:
        return {'status': 'fail', 'message': 'Entry doesn\'t exists'}

    prev = entry.prev_by_order()

    if not prev:
        return {'status': 'ok', 'message': 'First element'}

    prev.order, entry.order = entry.order, prev.order

    entry.save()
    prev.save()

    return {'status': 'ok'}


@app.route(
    '/api/category/move_down',
    endpoint='category.move_down',
    methods=['POST'],
)
@jsonify
def move_down():
    """
    Move category down
    """
    try:
        pk = int(request.form['pk'])
        entry = Category.get(Category.id == pk)
        """:type : Category """
    except ValueError:
        return {'status': 'fail', 'message': 'Wrong parameter'}
    except DoesNotExist:
        return {'status': 'fail', 'message': 'Entry doesn\'t exists'}

    prev = entry.next_by_order()

    if not prev:
        return {'status': 'ok', 'message': 'Last element'}

    prev.order, entry.order = entry.order, prev.order

    entry.save()
    prev.save()

    return {'status': 'ok'}
