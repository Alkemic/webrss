#!/usr/bin/env python
# -*- coding: utf-8 -*-
import os.path
import logging.config

from webrss.main import app
import webrss.views

if __name__ == "__main__":
    app.debug = True
    app.run(host='0.0.0.0')
