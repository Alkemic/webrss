# WebRSS

Web RSS client written in Python using Flask, peewee and AngularJS.

## Installation

* Clone this repo ``git clone https://github.com/Alkemic/webrss``
* Install libraries required for ``lxml``: ``aptitude install libxml2-dev libxslt1-dev python-dev``
* Install all dependencies ``pip install -r requirements.txt``
* Install database ``./scripts/create_tables.py``
* Add feed collector to crontab ``./scripts/collect_feeds.py``
* For testing use ``./run.py``

### Using MySQL/MariaDB

* Install client dev packages: ``sudo aptitude install libmysqlclient-dev`` or ``sudo aptitude install libmariadbclient-dev``
* Install mysql client ``pip install mysqlclient==1.3.10``, at this moment version 1.3.10 is fully working

## TODO

* Authorisation
* Unittest
