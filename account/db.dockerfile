FROM postgres:17.4
COPY up.sql /docker-entrypoint-initdb.d/1.sql

CMD ["postgres"]
