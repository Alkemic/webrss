# WebRSS

Web RSS client written in Python using Flask, peewee and AngularJS.

## Installation

* Clone this repo ``git clone https://github.com/Alkemic/webrss``
* Install libraries required for ``lxml``: ``aptitude install libxml2-dev libxslt1-dev python-dev``
* Install all dependencies ``pip install -r requirements.txt``
* Install database ``./scripts/create_tables.py``
* Add feed collector to crontab ``./scripts/collect_feeds.py``
* For testing use ``./run.py``

## TODO

* Authorisation
* Unittest
