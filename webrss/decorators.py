# -*- coding:utf-8 -*-
from functools import wraps

from flask import json

jsonify = lambda func: wraps(func)(lambda *func_args, **func_kwargs: json.jsonify(func(*func_args, **func_kwargs)))
