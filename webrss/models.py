# -*- coding:utf-8 -*-
"""
Models used in aplication
"""

from datetime import datetime

import peewee


DATABASE = peewee.SqliteDatabase('./webrss.db')


class BaseModel(peewee.Model):
    """ Base model class """

    class Meta:
        """ Base Meta model """
        database = DATABASE


class Category(BaseModel):
    """ Model containing all categories """
    title = peewee.CharField(max_length=255)

    order = peewee.IntegerField()

    created_at = peewee.DateTimeField(default=datetime.now())
    updated_at = peewee.DateTimeField(null=True)
    deleted_at = peewee.DateTimeField(null=True)

    class Meta:
        order_by = ('order',)

    def __unicode__(self):
        return self.title

    def __str__(self):
        return self.__unicode__()

    def not_deleted_feeds(self):
        """
        :rtype : list[Feed]
        """
        return self.feed_set.where(Feed.deleted_at.__eq__(None))

    def prev_by_order(self):
        """
        :rtype : Category
        """
        try:
            return Category.select() \
                .where(Category.order < self.order) \
                .where(Category.deleted_at.__eq__(None)) \
                .order_by(Category.order.desc()) \
                .limit(1)[0]
        except IndexError:
            return None

    def next_by_order(self):
        """
        :rtype : Category
        """
        try:
            return Category.select() \
                .where(Category.order > self.order) \
                .where(Category.deleted_at.__eq__(None)) \
                .order_by(Category.order.asc()) \
                .limit(1)[0]
        except IndexError:
            return None


class Feed(BaseModel):
    """ Model containing feeds """
    feed_title = peewee.CharField(max_length=255)
    feed_url = peewee.CharField(max_length=255)
    feed_image = peewee.CharField(max_length=255, null=True)
    feed_subtitle = peewee.CharField(max_length=255, null=True)

    site_url = peewee.CharField(max_length=255, null=True)
    site_favicon_url = peewee.CharField(max_length=255, null=True)

    category = peewee.ForeignKeyField(Category, null=True)

    created_at = peewee.DateTimeField(default=datetime.now())
    updated_at = peewee.DateTimeField(null=True)
    deleted_at = peewee.DateTimeField(null=True)

    class Meta:
        order_by = ('-feed_title',)

    def __unicode__(self):
        return self.feed_title

    def __str__(self):
        return self.__unicode__()

    def count_un_read(self):
        """
        Returns amount of unread entries
        """
        return self.entry_set.where(Entry.read_at.__eq__(None)).count()


class Entry(BaseModel):
    """ Model containing entries """
    title = peewee.CharField()
    author = peewee.CharField(null=True)
    summary = peewee.TextField(null=True)
    link = peewee.CharField()
    published_at = peewee.DateTimeField()
    feed = peewee.ForeignKeyField(Feed)

    read_at = peewee.DateTimeField(null=True)

    created_at = peewee.DateTimeField(default=datetime.now())
    updated_at = peewee.DateTimeField(null=True)
    deleted_at = peewee.DateTimeField(null=True)

    class Meta:
        order_by = ('-published_at',)

    def __unicode__(self):
        return self.title

    def __str__(self):
        return self.__unicode__()

    @property
    def is_read(self):
        """ Is this entry read? """
        return self.read_at is not None
