# ask-install

The `ask-install` tool can be used to stand up a group of Docker containers, simplifying the installation process for Coral Project's Ask product.

## Getting Started

- Option 1: Head over to the
[Latest release](https://github.com/coralproject/ask-install/releases/latest/)
page and download the precompiled binary that matched your environment. Currently the most popular is the linux option (ask-install_*_linux_amd64.tar.gz) since most newsrooms are setting up Ask on Amazon AWS via the EC2 service. . Run the binary and follow the instructions provided.

- Option 2: If you would rather you can also
compile from source provided you have a Go environment setup.


## System Requirements For The Ask Installer

The source system you use to build the Ask enviornment [needs the following installed](https://docs.coralproject.net/products/ask/#software-versions):

- [Docker - 1.12.1 or later](http://www.docker.com/products/docker/)
- [Docker Compose - 1.8.1 or later](http://www.docker.com/products/docker-compose/)

## Production Use

If this will be used in a production environment, it is strongly recommended you enable SSL and and link the instance to a hostname (example format using a subdomain mapped from your hosted instance to a DNS A record: https://your-subdomain.domain.com) for the machine. When you run the setup you must have the DNS records in place and working.
