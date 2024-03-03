#!/bin/sh

CREDENTIALS="host=$DBHOST user=$DBUSER password=$DBPASSWORD dbname=$DBNAME sslmode=$DBSSL"

cd migrations/

../goose postgres "$CREDENTIALS" up