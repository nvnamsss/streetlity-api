#/bin/bash

docker image rm -f api_container
docker build . -t api_container

