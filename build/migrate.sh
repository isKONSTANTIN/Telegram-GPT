#!/bin/sh

CREDENTIALS="host=$DB_HOST user=$DB_USER password=$DB_PASSWORD dbname=$DB_TABLE sslmode=$DBSSL"

cd migrations/

../goose postgres "$CREDENTIALS" up