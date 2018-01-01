#!/bin/bash

docker-compose rm -f && docker-compose up --build --force-recreate --abort-on-container-exit