#!/usr/bin/env python
# -*- coding: utf-8 -*-
"""
Run dev server
"""
from webrss.main import app
import webrss.views

if __name__ == "__main__":
    app.debug = True
    app.run(host='0.0.0.0', port=4567)
