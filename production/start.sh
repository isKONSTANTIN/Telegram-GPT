#!/bin/bash

RED='\033[0;31m'
NC='\033[0m'

if test -f "database_password"; then
  POSTGRES_PASSWORD=$(cat ./database_password)
else
  POSTGRES_PASSWORD=$(cat /dev/urandom | tr -dc '[:alpha:]' | fold -w ${1:-30} | head -n 1)
  echo $POSTGRES_PASSWORD > database_password
fi

POSTGRES_PASSWORD=$POSTGRES_PASSWORD docker-compose up -d

if [[ $? -eq 0 ]]; then
  echo -e "\nBot started.\nDocker Compose should automatically start containers after system reboot."
else
  echo -e "\n${RED}Start failed.$NC"
fi