#!/usr/bin/env bash

#TODO: 
# - create private network,
# - run hellosvc in network with no published ports
# - run gateway in network publishing port 443
#   and using volumes to give it access to your
#   cert and key files in `gateway/tls`
#   and setting environment variables
#    - CERTPATH = path to cert file in container
#    - KEYPATH = path to private key file in container
#    - HELLOADDR = net address of hellosvc container
if [ -z "$(docker network ls -q -f name=appnet)" ]
then
    docker network create appnet
fi 

docker run -d \
--name chatsvc \
--network appnet \
aethan/botsvc

docker run -d \
--name apiserver \
--network appnet \
-p 80:80 \
-p 443:443 \
-e REDISADDR=redissvr:6379 \
-e DBADDR=mongosvr:27017 \
-e BOTSVCADDR=botsvc \
-v /etc/letsencrypt:/etc/letsencrypt:ro \
-e TLSKEY=$TLSKEY \
-e TLSCERT=$TLSCERT \
-e EMAILPASS=$EMAILPASS \
-e SESSIONKEY=$SESSIONKEY \
aethan/apiserver