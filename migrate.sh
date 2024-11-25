#!/bin/bash

# Set the database URL (you can change this based on your environment)
DB_URL="postgres://postgres:1@localhost:5432/bcca_crawler"

# Function to run the "up" migrations
migrate_up() {
    echo "Running 'goose up' migration..."
    goose -dir "./sql/schema/" postgres $DB_URL up
}

# Function to run the "down" migrations
migrate_down() {
    echo "Running 'goose down' migration..."
    goose -dir "./sql/schema/" postgres $DB_URL down
}

# Check for command-line argument (either "up" or "down")
if [ "$1" == "up" ]; then
    migrate_up
elif [ "$1" == "down" ]; then
    migrate_down
else
    echo "Usage: $0 {up|down}"
    exit 1
fi
