#!/bin/bash

# Check if a file path was provided as an argument
if [ "$#" -ne 1 ]; then
    echo "Usage: $0 <secret>"
    exit 1
fi


echo "SECRET_KEY=$1" > .env
