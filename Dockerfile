FROM golang

MAINTAINER marco.gu

COPY ./backend ./
COPY ./templates/ ./templates/
COPY ./statics/ ./statics/

ENV OSIN_DB_USERNAME root
ENV OSIN_DB_PASSWORD abiglionnotfat
ENV OSIN_DB_HOST mysql-db
ENV OSIN_DB_PORT 3306
ENV OSIN_DB_DATABASE osintest

ENV WEB_DB_USERNAME root
ENV WEB_DB_PASSWORD abiglionnotfat
ENV WEB_DB_HOST mysql-db
ENV WEB_DB_PORT 3306
ENV WEB_DB_DATABASE dscenter

ENV GIN_MODE debug


ENTRYPOINT ./backend
