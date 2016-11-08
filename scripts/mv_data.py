#!/usr/bin/env python
# -*- coding: utf-8 -*-
import sys
sys.path.insert(0, '.')
from datetime import datetime

import peewee

from webrss import models


# OLD_DATABASE = peewee.SqliteDatabase(DB_NAME)
NEW_DATABASE = peewee.MySQLDatabase(
    'webrss',
    user='root',
    password='',
    host='127.0.0.1',
    port=13306,
    charset='utf8mb4',
)


class BaseModel(peewee.Model):
    """ Base model class """

    class Meta:
        """ Base Meta model """
        database = NEW_DATABASE


class Category(BaseModel):
    """ Model containing all categories """
    title = peewee.CharField(max_length=255)

    order = peewee.IntegerField()

    created_at = peewee.DateTimeField(default=datetime.now())
    updated_at = peewee.DateTimeField(null=True)
    deleted_at = peewee.DateTimeField(null=True)


class Feed(BaseModel):
    """ Model containing feeds """
    feed_title = peewee.CharField(max_length=255)
    feed_url = peewee.CharField(max_length=255)
    feed_image = peewee.CharField(max_length=255, null=True)
    feed_subtitle = peewee.TextField(null=True)

    site_url = peewee.CharField(max_length=255, null=True)
    site_favicon_url = peewee.CharField(max_length=255, null=True)

    category = peewee.ForeignKeyField(Category, null=True)

    last_read_at = peewee.DateTimeField(default=datetime.now())
    created_at = peewee.DateTimeField(default=datetime.now())
    updated_at = peewee.DateTimeField(null=True)
    deleted_at = peewee.DateTimeField(null=True)


class Entry(BaseModel):
    """ Model containing entries """
    title = peewee.CharField(max_length=512)
    author = peewee.CharField(null=True)
    summary = peewee.TextField(null=True)
    link = peewee.CharField()
    published_at = peewee.DateTimeField()
    feed = peewee.ForeignKeyField(Feed)

    read_at = peewee.DateTimeField(null=True)

    created_at = peewee.DateTimeField(default=datetime.now())
    updated_at = peewee.DateTimeField(null=True)
    deleted_at = peewee.DateTimeField(null=True)


MODELS = (
    (models.Category, Category),
    (models.Feed, Feed),
    (models.Entry, Entry),
)
for old_cls, cls in MODELS:
    for obj in old_cls.select():
        try:
            cls.create(**obj._data).save()
        except Exception as e:
            print e
