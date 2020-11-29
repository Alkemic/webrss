FROM busybox

FROM golang:1.15 as backend

COPY . /build
WORKDIR /build
RUN go mod download && \
    CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -o webrss-app ./cmd/webrss && \
    strip webrss-app
ADD https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz migrate.linux-amd64.tar.gz
RUN tar zxvf migrate.linux-amd64.tar.gz

FROM node:6 as frontend

COPY . /build
WORKDIR /build
RUN cd frontend && \
    npm i && \
    PRODUCTION=true ./node_modules/.bin/gulp build

FROM scratch

COPY --from=backend /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=backend /build/migrations /migrations
COPY --from=backend /build/webrss-app /webrss
COPY --from=backend /build/migrate.linux-amd64 /migrate
COPY --from=frontend /build/static /static/
COPY --from=busybox /bin/busybox /bin/busybox
COPY --from=busybox /bin/sh /bin/sh
COPY templates/* /templates/

CMD /migrate -path /migrations/ -database "mysql://${DB_DSN}" up; /webrss
