# -*- coding:utf-8 -*-
"""
Decorators
"""
from functools import wraps

from flask import json


def jsonify(func):
    """ Returns data as JSON """

    @wraps(func)
    def wrapped(*func_args, **func_kwargs):
        """The inner function """
        return json.jsonify(func(*func_args, **func_kwargs))

    return wrapped
