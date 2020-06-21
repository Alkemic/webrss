# WebRSS

Web RSS client written in Go and AngularJS.

## Running (dockerised)
* ``docker run --name webrss -d alkemic/webrss``
* See [docker-compose.yml](./docker-compose.yml) for details on usage with docker compose

## Running

* Build backend
  * Get source ``git clone https://github.com/Alkemic/webrss/``
  * Build & install ``go install ./cmd/webrss/``
* Build frontend
  * Download dependencies ``npm i``
  * Build frontend ``./node_modules/.bin/gulp build``
* Setup environment varaibles
  * `DB_DSN` - the DSN for database, ie: `root:@tcp(127.0.0.1:13306)/webrss?parseTime=true`, mind the `parseTime=true` part
  * `BIND_ADDR` - bind address, ie: `:8080`
  * (optional) `PER_PAGE` - how many entries will be loaded when feed is selected
* Run from main folder ``webrss``

## Database

* Install [golang migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate#installation)
* Create database ``CREATE DATABASE `webrss3` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;``, 
utf8mb4 is required, as many sites uses emojis
* Migrate using ``migrate -path ./migrations/ -database "mysql://root:toor@tcp(localhost:3306)/webrss" up``
