#!/bin/bash -u
docker-compose down
docker rm $(docker ps -aq)
rm -rf fixtures/channel-artifacts fixtures/crypto-config
docker volume prune
