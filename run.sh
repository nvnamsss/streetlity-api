#/bin/bash
name=$1

echo "STOP RUNNING API CONTAINER"
docker stop -t 30 ${name}_api_container 
docker rm -f ${name}_api_container

echo "DONE STOPPING"

docker run --name ${name}_api_container -d\
            --network common-net \
            --restart always \
            -p 9000:9000 \
            api_container

    
