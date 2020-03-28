FROM golang:1.14 as backend

COPY . /build
WORKDIR /build
RUN go mod download && \
    CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -o webrss-app ./cmd/webrss && \
    strip webrss-app

FROM node:6 as frontend

COPY . /build
WORKDIR /build
RUN cd frontend && \
    npm i && \
    PRODUCTION=true ./node_modules/.bin/gulp build

FROM scratch

COPY --from=backend /build/webrss-app /webrss
COPY --from=frontend /build/static/* /static/
COPY templates/* /templates/

ENTRYPOINT ["/webrss"]
