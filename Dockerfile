FROM debian:9

ENV APP_DIR /app

ENV PACKAGES_REQUIRED_PRE_BUILD gnupg2 apt-transport-https ca-certificates

ENV PACKAGES_REQUIRED_BUILD libxml2-dev libxslt1-dev \
    lsb-release \
    python-dev python-pip python-setuptools \
    curl gnupg2 ca-certificates gcc nodejs
ENV PACKAGES_REQUIRED python libmariadbclient-dev-compat libmariadbclient18 \
    libxml2 libpython2.7 libxslt1.1

ENV NODE_SOURCES_LIST "deb [arch=amd64] https://deb.nodesource.com/node_6.x stretch main"

COPY . $APP_DIR

WORKDIR $APP_DIR

ADD https://deb.nodesource.com/gpgkey/nodesource.gpg.key /tmp/nodesource.gpg.key

RUN set -x && \
    apt-get update && \
    apt-get install -y --no-install-recommends $PACKAGES_REQUIRED_PRE_BUILD && \
    \
    echo $NODE_SOURCES_LIST > /etc/apt/sources.list.d/nodesource.list && \
    cat /etc/apt/sources.list.d/nodesource.list && \
    apt-key add /tmp/nodesource.gpg.key && \
    \
    apt-get update && \
    apt-get install -y --no-install-recommends $PACKAGES_REQUIRED $PACKAGES_REQUIRED_BUILD && \
    \
    pip install -r requirements.txt && \
    pip install uwsgi==2.0.15 && \
    (cd frontend && npm i && ./node_modules/.bin/gulp build && rm -rf ./node_modules/) && \
    \
    apt-get remove --purge -y $PACKAGES_REQUIRED_PRE_BUILD $PACKAGES_REQUIRED_BUILD && \
    apt-get install -y --no-install-recommends ca-certificates&& \
    apt-get clean && \
    apt-get autoremove -y && \
    apt-get autoclean -y && \
    rm -rf /root/.npm/* \
        /root/.cache/pip/* \
        /var/lib/apt/lists/* /tmp/* /var/tmp/*

EXPOSE 8000

CMD ./scripts/create_tables.py && uwsgi -y $APP_DIR/uwsgi.yaml
