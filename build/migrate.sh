#!/bin/sh

CREDENTIALS="host=$DB_HOST user=$DB_USER password=$DB_PASSWORD dbname=$DB_TABLE sslmode=$DB_SSL"

cd migrations/

../goose postgres "$CREDENTIALS" up