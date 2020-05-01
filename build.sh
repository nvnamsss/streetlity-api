#/bin/bash

docker image rm user_container
docker build . -t user_container

