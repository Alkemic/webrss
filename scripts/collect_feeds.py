#!/usr/bin/env python
#  -*- coding:utf-8 -*-
import sys

sys.path.insert(0, '.')  # noqa

import threading
import requests

from webrss.functions import process_feed
from webrss.models import Feed, Category


threads = []


class FeedThread(threading.Thread):
    def __init__(self, feed):
        threading.Thread.__init__(self)
        self.feed = feed

    def run(self):
        print u"started, {}".format(self.feed.__str__())
        try:
            resp = requests.get(self.feed.feed_url, timeout=30)
            feed_data = resp.content
            process_feed(self.feed, feed_data)
        except (requests.ConnectionError, requests.ReadTimeout) as e:
            print "exception at", self.feed, e


feeds = Feed.select().where(
    Feed.deleted_at >> None,
    Category.deleted_at >> None,
).join(
    Category,
)
for feed in feeds:
    t = FeedThread(feed)
    t.start()
    threads.append(t)

for t in threads:
    t.join()
