#!/bin/bash

# Check if an argument is provided
if [ "$#" -ne 1 ]; then
    echo "Usage: $0 <sql-connection-string>"
    exit 1
fi

# Store the SQL connection string from the input argument
SQL_CONNECTION_STRING="$1"

# Run your Go programs with the provided SQL connection string
go run microservices/auth/*.go -sql "$SQL_CONNECTION_STRING" &
go run microservices/users/*.go -sql "$SQL_CONNECTION_STRING" &
go run microservices/records/*.go -sql "$SQL_CONNECTION_STRING" &
