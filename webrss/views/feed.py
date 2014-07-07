# -*- coding:utf-8 -*-
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
    entries = [entry for i, entry in enumerate(feed.entry_set.where(Entry.deleted_at == None))]

    return render_template('feed/index.html', entries=entries)


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

    try:
        feed_kwargs = {'feed_url': request.form['url'], 'category': category}
        feed = feedparser.parse(request.form['url'])

        feed_kwargs['feed_title'] = feed.feed['title']
        feed_kwargs['feed_image'] = feed.feed['image'] if 'image' in feed.feed else None
        feed_kwargs['feed_subtitle'] = feed.feed['subtitle'] if 'subtitle' in feed.feed else None

        if 'link' in feed.feed:
            feed_kwargs['site_url'] = feed.feed['link']
            feed_kwargs['site_favicon_url'] = get_favicon(feed.feed['link'])

        entry = Feed.create(**feed_kwargs)
        process_feed(entry)
        return {'status': 'ok'}
    except Exception as e:  # todo: don't catch'em all.
        return {'status': 'fail'}


@app.route('/api/feed/update/<int:pk>', endpoint='feed.update', methods=['POST'])
@jsonify
def update(pk):
    """
    Update content of feed
    todo:
    """
    entry = Feed.get(Feed.id == pk)

    return {}


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
        return {'status': 'fail', 'message': 'Exception occurred during deleting'}

    return {'status': 'ok'}
