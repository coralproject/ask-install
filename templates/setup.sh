#!/bin/sh

set -e

echo "Pulling service layers:"
echo ""

docker-compose pull

echo "Layers pulled."
echo ""

echo "Creating the services now:"

echo ""

# Start the services.
docker-compose up -d

echo ""
echo "Services created."
{{ if .Password }}
echo ""
echo "Now creating user:"
echo ""

# Create the user.
docker-compose run --rm auth cli create -f --email "{{ .Email }}" --password "{{ .Password }}" --name "{{ .DisplayName }}"

echo ""
echo "User was created."
echo ""
echo "Please remove this file as it contains the plaintext password."
{{ end }}
