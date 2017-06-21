#!/usr/bin/env python
#  -*- coding:utf-8 -*-
import sys

sys.path.insert(0, '.')

import threading
import requests

from webrss.functions import process_feed
from webrss.models import Feed


threads = []


class FeedThread(threading.Thread):
    def __init__(self, feed):
      threading.Thread.__init__(self)
      self.feed = feed

    def run(self):
        try:
            feed_data = requests.get(self.feed.feed_url).content
            process_feed(self.feed, feed_data)
        except requests.ConnectionError:
            pass


for feed in Feed.select():
    t = FeedThread(feed)
    t.start()
    threads.append(t)

for t in threads:
    t.join()
