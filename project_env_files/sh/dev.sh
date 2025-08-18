#!/bin/bash

# port="$(cat .env-DEV | grep -m1 'PORT' | cut -c 8-)"
# echo $port
# echo "PORT = $port" >> .env


# export ENVIRONMENT=DEV
# go run .

# Set environment to DEV
export ENVIRONMENT=DEV

# Optional: print which environment is used
echo "Running in environment: $ENVIRONMENT"

# Optional: load PORT from .env-DEV without modifying .env
if [ -f ".env-DEV" ]; then
    export PORT=$(grep -m1 '^PORT=' .env-DEV | cut -d '=' -f2)
    echo "Using PORT: $PORT"
fi

# Run the Go application
go run .

