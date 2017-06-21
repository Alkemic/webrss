# -*- coding:utf-8 -*-
"""
Main WebRSS module
"""
import os

import peewee


DB_TYPE = os.environ.get('DB_TYPE', 'sqlite')
DB_HOST = os.environ.get('DB_HOST')
DB_PORT = os.environ.get('DB_PORT')
DB_NAME = os.environ.get('DB_NAME', './webrss.db')
DB_USER = os.environ.get('DB_USER')
DB_PASS = os.environ.get('DB_PASS')

if DB_TYPE not in ('mysql', 'sqlite'):
    raise ValueError('`DB_TYPE` must be \'mysql\' or \'sqlite\'')

if DB_TYPE == 'sqlite':
    DATABASE = peewee.SqliteDatabase(DB_NAME, threadlocals=True)
elif DB_TYPE == 'mysql':
    DATABASE = peewee.MySQLDatabase(
        DB_NAME,
        user=DB_USER,
        password=DB_PASS or '',
        host=DB_HOST or '127.0.0.1',
        port=int(DB_PORT or 3306),
        charset='utf8mb4',
        threadlocals=True,
    )
