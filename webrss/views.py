# -*- coding:utf-8 -*-
"""Feed related views"""
from datetime import datetime

import feedparser
from flask import render_template
from flask import request
from flask_peewee.rest import RestResource
import peewee

from .models import Feed, Category, Entry
from .main import app, rest_api
from .decorators import jsonify
from .functions import get_favicon, process_feed

RestResource.paginate_by = 50
RestResource.authorize = lambda *args, **kwargs: True


class CategoryResource(RestResource):
    exclude = ('created_at', 'updated_at', 'deleted_at',)

    def prepare_data(self, obj, data):
        data['feeds'] = [
            {
                'id': feed.id,
                'feed_title': feed.feed_title,
                'feed_image': feed.feed_image,
                'feed_url': feed.feed_url,
                'site_favicon_url': feed.site_favicon_url,
                'category': feed.category.id,
                'last_read_at': str(feed.last_read_at),
                'un_read': feed.count_un_read(),
                'new_entries': feed.last_read_at < feed.last_entry.created_at,
            }
            for feed in obj.not_deleted_feeds()
        ]
        return data


class FeedResource(RestResource):
    exclude = ('created_at', 'updated_at', 'deleted_at',)
    include_resources = {'category': CategoryResource}

    def prepare_data(self, obj, data):
        data['un_read'] = obj.count_un_read()
        return data

    def save_object(self, instance, raw_data):
        if self.pk:
            return super(FeedResource, self).save_object(instance, raw_data)

        try:
            instance.category = Category\
                .get(Category.id == raw_data['category'])
        except peewee.DoesNotExist:
            pass

        feed = feedparser.parse(instance.feed_url)
        instance.feed_title = feed.feed['title']

        if 'image' in feed.feed:
            instance.feed_image = feed.feed['image']

        if 'subtitle' in feed.feed:
            instance.feed_subtitle = feed.feed['subtitle']

        if 'link' in feed.feed:
            instance.site_url = feed.feed['link']
            instance.site_favicon_url = get_favicon(feed.feed['link'])

        process_feed(instance)

        return super(FeedResource, self).save_object(instance, raw_data)


class EntryResource(RestResource):
    exclude = ('created_at', 'updated_at', 'deleted_at',)
    include_resources = {'feed': FeedResource}

    def prepare_data(self, obj, data):
        data['new_entry'] = obj.created_at > obj.feed.last_read_at
        return super(EntryResource, self).prepare_data(obj, data)

    def object_list(self):
        object_list = super(EntryResource, self).object_list()

        if 'feed' in request.args:
            feed = Feed.get(id=request.args['feed'])
            feed.last_read_at = datetime.now()
            feed.save()

        return object_list

    def object_detail(self, obj):
        obj.read_at = datetime.now()
        obj.save()

        return super(EntryResource, self).object_detail(obj)


rest_api.register(Feed, FeedResource)
rest_api.register(Category, CategoryResource)
rest_api.register(Entry, EntryResource)

rest_api.setup()


@app.route('/')
def index():
    """Index view"""
    return render_template('index.html')


@jsonify
@app.route('/api/search', methods=['POST'])
def search():
    """Return search result"""
    phrase = '%%%s%%' % request.form['phrase']

    entries = Entry.select()\
        .where((Entry.title ** phrase) | (Entry.summary ** phrase))\
        .where(Entry.deleted_at.__eq__(None))

    return render_template('search.html', entries=list(entries))


@app.route('/api/category/<int:pk>/move_up', methods=['POST'])
@jsonify
def move_up(pk):
    """
    Move category up
    """
    try:
        entry = Category.get(Category.id == pk)
        """:type : Category """
    except ValueError:
        return {'status': 'fail', 'message': 'Wrong parameter'}
    except peewee.DoesNotExist:
        return {'status': 'fail', 'message': 'Entry doesn\'t exists'}

    prev = entry.prev_by_order()

    if not prev:
        return {'status': 'ok', 'message': 'First element'}

    prev.order, entry.order = entry.order, prev.order

    try:
        entry.save()
        prev.save()

        return {'status': 'ok'}
    except peewee.DatabaseError:
        return {
            'status': 'fail',
            'message': 'Exception occurred during saving'
        }


@app.route('/api/category/<int:pk>/move_down', methods=['POST'])
@jsonify
def move_down(pk):
    """
    Move category down
    """
    try:
        entry = Category.get(Category.id == pk)
        """:type : Category """
    except ValueError:
        return {'status': 'fail', 'message': 'Wrong parameter'}
    except peewee.DoesNotExist:
        return {'status': 'fail', 'message': 'Entry doesn\'t exists'}

    prev = entry.next_by_order()

    if not prev:
        return {'status': 'ok', 'message': 'Last element'}

    prev.order, entry.order = entry.order, prev.order

    try:
        entry.save()
        prev.save()

        return {'status': 'ok'}
    except peewee.DatabaseError:
        return {
            'status': 'fail',
            'message': 'Exception occurred during saving'
        }
