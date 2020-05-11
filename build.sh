ui#/bin/bash

docker image rm api_container
docker build . -t api_container

