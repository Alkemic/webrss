#!/usr/bin/env python
# -*- coding: utf-8 -*-
import sys
sys.path.insert(0, '.')

from webrss import models


for cls in models.Category, models.Feed, models.Entry:
    cls.create_table()
