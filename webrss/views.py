# -*- coding:utf-8 -*-
"""Feed related views"""
import base64
from datetime import datetime
from urllib import urlencode

import feedparser
from flask import render_template
from flask import request
from flask_peewee.rest import RestResource
import peewee
from playhouse.shortcuts import model_to_dict

from functions import get_favicon_url
from . import DATABASE
from .models import Feed, Category, Entry
from .main import app, rest_api
from .decorators import jsonify
from .functions import get_favicon, process_feed, to_datetime

PER_PAGE = 50
ENTRY_LAST_PUBLISHED_AT_SQL = """
select
    f.id,
    (
        select max(e.published_at)
        from entry e
        where e.feed_id = f.id
    ) entry_max_published_at
from feed f
where f.deleted_at is null;
"""


RestResource.paginate_by = PER_PAGE
RestResource.authorize = lambda *args, **kwargs: True


class CategoryResource(RestResource):
    _entry_last_published_at = None
    _entry_last_published = None

    exclude = ('created_at', 'updated_at', 'deleted_at',)

    @property
    def entry_last_published_at(self):
        invalidated = (
            not self._entry_last_published or
            (datetime.now() - self._entry_last_published_at).seconds > 5
        )
        if invalidated:
            cursor = DATABASE.execute_sql(ENTRY_LAST_PUBLISHED_AT_SQL)
            self._entry_last_published = {
                pk: to_datetime(published_ad) if published_ad else None
                for pk, published_ad in cursor.fetchall()
            }
            self._entry_last_published_at = datetime.now()

        return self._entry_last_published

    def prepare_data(self, obj, data):
        data['feeds'] = []
        for feed in obj.not_deleted_feeds():
            data['feeds'].append({
                'id': feed.id,
                'feed_title': feed.feed_title,
                'feed_image': feed.feed_image,
                'feed_url': feed.feed_url,
                'site_favicon_url': feed.site_favicon_url,
                'site_favicon': feed.site_favicon,
                'category': feed.category.id,
                'last_read_at': str(feed.last_read_at),
                'un_read': feed.count_un_read,
                'new_entries': feed.has_new_entries(
                    self.entry_last_published_at[feed.id],
                ),
            })

        return data


class FeedResource(RestResource):
    exclude = ('created_at', 'updated_at', 'deleted_at',)
    include_resources = {'category': CategoryResource}

    def prepare_data(self, obj, data):
        data['un_read'] = obj.count_un_read
        return data

    def save_object(self, instance, raw_data):
        if instance.id:
            favicon = get_favicon(instance.site_favicon_url)
            instance.site_favicon = (
                base64.b64encode(favicon) if favicon else None
            )
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

            if not instance.site_favicon_url:
                favicon_url = get_favicon_url(feed.feed['link'])
                instance.site_favicon_url = favicon_url

        favicon = get_favicon(instance.site_favicon_url)
        instance.site_favicon = base64.b64encode(favicon) if favicon else None

        returned = super(FeedResource, self).save_object(instance, raw_data)

        process_feed(instance)

        return returned


class EntryResource(RestResource):
    exclude = ('created_at', 'updated_at', 'deleted_at',)
    include_resources = {'feed': FeedResource}

    def prepare_data(self, obj, data):
        data['new_entry'] = obj.created_at > obj.feed.last_read_at
        data['published_at'] = obj.published_at.strftime("%Y-%m-%d %H:%M")
        return super(EntryResource, self).prepare_data(obj, data)

    def get_query(self):
        return self.model.select().join(Feed).where(
            Feed.deleted_at.is_null(True)
        )

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


@app.route("/api/search", methods=["GET"])
@jsonify
def search():
    phrase = request.args["phrase"]
    if not phrase:
        return []

    phrase = "%{}%".format(phrase)
    try:
        page = int(request.args.get("page", 1))
        page = page if page > 0 else 1
    except ValueError:
        page = 1

    categories = Category.select().where(Category.deleted_at.is_null(True))
    entries = Entry.select()\
        .join(Feed)\
        .where(Feed.category.in_(categories))\
        .where(Feed.deleted_at.is_null(True))\
        .where((Entry.title ** phrase) | (Entry.summary ** phrase))\
        .where(Entry.deleted_at.is_null(True))\
        .offset((page-1)*PER_PAGE)\
        .limit(PER_PAGE)

    query_params = request.args.to_dict()
    query_params.update({"page": page+1})
    data = {
        "objects": [model_to_dict(entry) for entry in entries],
        "meta": {},
    }
    print len(data["objects"])
    if len(data["objects"]) == PER_PAGE:
        data["meta"]["next"] = "/api/search?{}".format(urlencode(query_params))
    return data

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
