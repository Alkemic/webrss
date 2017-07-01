#!/usr/bin/env python
# -*- coding: utf-8 -*-
import sys
sys.path.insert(0, '.')

import peewee
from playhouse.migrate import SchemaMigrator, migrate

from webrss import models


for cls in models.Category, models.Feed, models.Entry:
    try:
        cls.create_table()
    except peewee.OperationalError as ex:
        print ex

migrator = SchemaMigrator(models.DATABASE)

try:
    try:
        migrate(migrator.drop_index('entry_feed_id_published_at'))
    except:
        pass

    migrate(migrator.add_index('entry', ('feed_id', 'published_at'), False))
except peewee.OperationalError as ex:
    print ex
