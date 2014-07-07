#!/usr/bin/env python
# -*- coding: utf-8 -*-
import sys
sys.path.insert(0, '.')

from webrss.main import app
from webrss import models


for cls in dir(models):
    cls = getattr(models, cls)
    is_table_class = all((
        isinstance(cls, type) and issubclass(cls, models.BaseModel),
        hasattr(cls, '__name__') and cls.__name__ != 'BaseModel',
    ))

    if is_table_class:
        cls.create_table()
