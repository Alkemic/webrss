#!/usr/bin/env python
#  -*- coding:utf-8 -*-
import sys

sys.path.insert(0, '.')

from webrss.functions import process_feed

from webrss.models import Feed


for feed in Feed.select():
    process_feed(feed)
