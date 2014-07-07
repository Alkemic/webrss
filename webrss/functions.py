# -*- coding:utf-8 -*-
from datetime import datetime
import sys
import urllib2
from urlparse import urlparse
import feedparser

import lxml.html
import peewee

from webrss.models import Category, Entry


def categories_dict():
    """
    Returns not deleted categories in a dict
    """
    entries = Category.select().where(Category.deleted_at == None).dicts()

    return {i: entry for i, entry in enumerate(entries)}


def categories_tuple():
    """
    Returns not deleted categories in a tuple
    """
    entries = Category.select().where(Category.deleted_at == None).dicts()

    return tuple(entry for i, entry in enumerate(entries))


def get_favicon(url):
    """
    Fetch favicon from given url
    First try, to find a link tag, then tries to fetch <domain>/favicon.ico

    :type url: str
    :return: str|bool
    """
    headers = {
        'User-Agent': 'urllib2 (Python %s)' % sys.version.split()[0],
        'Connection': 'close',
    }

    parsed = urlparse(url)
    url = '%s://%s/' % (parsed.scheme, parsed.netloc)

    request = urllib2.Request(url, headers=headers)
    try:
        content = urllib2.urlopen(request).read()
        icon_path = lxml.html.fromstring(content).xpath('//link[@rel="icon" or @rel="shortcut icon"]/@href')
        if icon_path:
            icon_path = icon_path[-1]
            if icon_path[:6] in ('http:/', 'https:', 'ftp://'):
                return icon_path
            else:
                return url + icon_path
    except(urllib2.HTTPError, urllib2.URLError) as e:
        pass

    request = urllib2.Request(url + 'favicon.ico', headers=headers)
    try:
        urllib2.urlopen(request).read()
        return url + 'favicon.ico'
    except(urllib2.HTTPError, urllib2.URLError):
        pass

    return None


def process_feed(feed):
    """
    For given feed
    :type feed: webrss.models.Feed
    """
    parsed = feedparser.parse(feed.feed_url)

    feed.feed_title = parsed.feed['title']
    feed.feed_image = parsed.feed['image'] if 'image' in parsed.feed else None
    feed.feed_subtitle = parsed.feed['subtitle'] if 'subtitle' in parsed.feed else None

    if 'link' in parsed.feed:
        feed.site_url = parsed.feed['link']
        feed.site_favicon_url = get_favicon(parsed.feed['link'])

    feed.save()

    for entry in parsed.entries:
        try:
            feed_entry = Entry.get(Entry.link == entry['link'])
        except peewee.DoesNotExist:
            feed_entry = Entry(link=entry['link'])

        feed_entry.title = entry['title']
        feed_entry.author = entry['author'] if 'author' in entry else None
        feed_entry.summary = entry['summary'] if 'summary' in entry else None
        feed_entry.link = entry['link']

        if 'published_parsed' in entry:
            feed_entry.published_at = datetime(*entry['published_parsed'][:6])
        elif 'updated_parsed' in entry:
            feed_entry.published_at = datetime(*entry['updated_parsed'][:6])
        else:
            feed_entry.published_at = datetime.now()

        feed_entry.feed = feed
        feed_entry.save()
