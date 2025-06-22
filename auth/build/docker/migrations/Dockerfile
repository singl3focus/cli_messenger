ARG MIGRATION_DIR

FROM alpine:3.18

ARG MIGRATION_DIR

RUN apk update && \
    apk upgrade && \
    apk add --no-cache bash tzdata && \
    cp /usr/share/zoneinfo/Europe/Moscow /etc/localtime && \
    echo "Europe/Moscow" > /etc/timezone && \
    apk del tzdata

ADD https://github.com/pressly/goose/releases/download/v3.14.0/goose_linux_x86_64 /bin/goose
RUN chmod +x /bin/goose

WORKDIR /opt/migrations

RUN mkdir -p migrations

COPY ${MIGRATION_DIR} ./migrations/

RUN echo "migration_dir=${MIGRATION_DIR}"
RUN ls -laR > /ls_output.txt && cat /ls_output.txt

CMD ["sh", "-c", "sleep 2 && goose -dir ./migrations postgres \"${PG_DSN}\" up -v"]
