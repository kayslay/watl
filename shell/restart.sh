#!/bin/bash

cd watl/
git pull origin master
docker-compose up -d --build watl