# -*- coding:utf-8 -*-
"""
Feed related views
"""

from flask import render_template
import peewee
import feedparser
from flask import request
from datetime import datetime

from webrss.decorators import jsonify
from webrss.functions import get_favicon, process_feed
from webrss.main import app
from webrss.models import Feed, Category, Entry


@app.route('/api/feed/index', endpoint='feed.index', methods=['POST'])
def index():
    """
    Returns list of not deleted categories
    """
    feed = Feed.get(Feed.id == request.form['feed_id'])
    feeds = list(feed.entry_set.where(Entry.deleted_at.__eq__(None)))

    return render_template('feed/index.html', entries=feeds)


@app.route('/api/feed/create', endpoint='feed.create', methods=['POST'])
@jsonify
def create():
    """
    Create new category
    """
    try:
        category = Category.get(Category.id == request.form['category'])
    except peewee.DoesNotExist:
        return {'status': 'fail', 'message': 'Category doesn\'t exists'}

    feed_kwargs = {'feed_url': request.form['url'], 'category': category}
    feed = feedparser.parse(request.form['url'])
    feed_kwargs['feed_title'] = feed.feed['title']

    if 'image' in feed.feed:
        feed_kwargs['feed_image'] = feed.feed['image']

    if 'subtitle' in feed.feed:
        feed_kwargs['feed_subtitle'] = feed.feed['subtitle']

    if 'link' in feed.feed:
        feed_kwargs['site_url'] = feed.feed['link']
        feed_kwargs['site_favicon_url'] = get_favicon(feed.feed['link'])

    try:
        entry = Feed.create(**feed_kwargs)
        process_feed(entry)
        return {'status': 'ok'}
    except peewee.DatabaseError:
        return {'status': 'fail'}


@app.route('/api/feed/update', endpoint='feed.update', methods=['POST', 'GET'])
@jsonify
def update():
    """
    Update content of feed
    """
    is_get = request.method == 'GET'

    if is_get:
        pk = request.args.get('pk', None)
    else:
        pk = request.form.get('pk', None)

    if pk is None:
        return {'status': 'fail', 'message': 'Wrong parameter'}

    try:
        entry = Feed.get(Feed.id == pk)
        """ :type : Feed """
    except ValueError:
        return {'status': 'fail', 'message': 'Wrong parameter'}
    except peewee.DoesNotExist:
        return {'status': 'fail', 'message': 'Entry doesn\'t exists'}

    if is_get:  # if requesting via GET, return feed data
        return {'category': entry.category.id, 'url': entry.feed_url}

    if not all(name in request.form for name in ('url', 'pk', 'category')):
        return {'status': 'fail', 'message': 'Not all parameters were sent'}

    try:
        category = Category.get(Category.id == request.form['category'])
    except ValueError:
        return {'status': 'fail', 'message': 'Wrong parameter'}
    except peewee.DoesNotExist:
        return {'status': 'fail', 'message': 'Category doesn\'t exists'}

    try:
        entry.feed_url = request.form['category']
        entry.category = category
        entry.save()
        return {'status': 'ok'}
    except peewee.DatabaseError:
        return {'status': 'fail', 'message': 'Exception occurred during save'}


@app.route('/api/feed/delete/', endpoint='feed.delete', methods=['POST'])
@jsonify
def delete():
    """
    Delete given feed, requires pk in POST table
    """
    try:
        pk = request.form['pk']
        entry = Feed.get(Feed.id == pk)
    except ValueError:
        return {'status': 'fail', 'message': 'Wrong parameter'}
    except peewee.DoesNotExist:
        return {'status': 'fail', 'message': 'Entry doesn\'t exists'}

    entry.deleted_at = datetime.now()
    try:
        entry.save()
    except peewee.DatabaseError:
        return {
            'status': 'fail',
            'message': 'Exception occurred during deleting'
        }

    return {'status': 'ok'}
