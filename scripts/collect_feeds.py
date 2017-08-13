#!/usr/bin/env python
#  -*- coding:utf-8 -*-
import socket
import sys
import threading

import httplib2

sys.path.insert(0, '.')  # noqa

from webrss.functions import process_feed
from webrss.models import Feed, Category


TIMEOUT = 45
HEADERS = {
    "user-agent": "WebRSS parser (https://github.com/Alkemic/webrss)",
}

threads = []


h = httplib2.Http(timeout=TIMEOUT)


class FeedThread(threading.Thread):
    def __init__(self, feed):
        threading.Thread.__init__(self)
        self.feed = feed

    def run(self):
        try:
            resp, content = h.request(
                self.feed.feed_url,
                headers=HEADERS,
            )
            process_feed(self.feed, content)
        except socket.timeout:
            print "Feed '{}' timeouted".format(self.feed.__str__())


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
