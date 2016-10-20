# ask-install

The `ask-install` tool can be used to bootstrap an environment to run the
Coral Project's Ask product.

## Getting Started

Head over to the
[Latest release](https://github.com/coralproject/ask-install/releases/latest)
page and download the precompiled binary there. If you would rather you can also
compile from source provided you have a Go environment setup.

Then just run the binary and follow the instructions provided.

## System Requirements

- [Docker](https://www.docker.com/)
- [Docker Compose](https://www.docker.com/products/docker-compose)

## Production Use

If this is to be used in production, it is required that you enable SSL and
provide a real hostname for the machine. When you run this setup you must have
already pointed your DNS records to your machine running these services.

This installer uses the [Caddy Webserver](https://caddyserver.com/) which will
automatically setup and manage SSL for you provided that the DNS records point
to the machine and the docker user can bind to ports 80/443.
