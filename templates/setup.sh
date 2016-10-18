#!/bin/sh

set -e

echo "Creating the services now:"

echo ""

# Start the services.
docker-compose up -d

echo ""
echo "Services created. Now creating user:"
echo ""

# Create the user.
docker-compose run --rm auth cli create -f --email "{{ .Email }}" --password "{{ .Password }}" --name "{{ .DisplayName }}"

echo ""
echo "User was created."
echo ""
echo "Done. Please remove this file as it contains the plaintext password."
