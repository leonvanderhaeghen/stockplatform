#!/bin/sh

# Set the MongoDB server hostname
export ME_CONFIG_MONGODB_SERVER=stockplatform-mongodb

# Wait for MongoDB to be ready
echo "Waiting for MongoDB to be ready..."
until nc -z stockplatform-mongodb 27017; do
  echo "Waiting for MongoDB..."
  sleep 1
done

# Start MongoDB Express
echo "Starting MongoDB Express..."
exec node app.js
