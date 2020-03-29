#/bin/bash
name=$1

echo "STOP RUNNING API CONTAINER"
docker stop -t 30 {$1}_api_container 
docker rm -f {$1}_api_container

echo "DONE STOPPING"

docker run --name {$1}_api_container \
            --network common_net \
            --restart always \
            -p 9000:9000 \
            -d api_container:latest

    