FROM postgres:16.3-alpine

COPY docker/postgres/init.sql /docker-entrypoint-initdb.d/01-init.sql

EXPOSE 5432