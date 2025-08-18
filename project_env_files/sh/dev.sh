#!/bin/bash

# port="$(cat .env-DEV | grep -m1 'PORT' | cut -c 8-)"
# echo $port
# echo "PORT = $port" >> .env

export ENVIRONMENT=DEV
go run .
