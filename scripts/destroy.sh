#!/bin/bash

# this program lists all of the droplets behind a load balancer
# and asks a user to destroy one of them

get_balancer()
{
    resp=$(curl -X GET \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" "https://api.digitalocean.com/v2/load_balancers")
    balancer=$(echo $resp | jq -r '.load_balancers[0].id')
    echo $balancer
    droplet_ids=$(echo $resp | jq '.load_balancers[0].droplet_ids')
    echo $droplet_ids
}

delete_droplet_from_balancer()
{
    echo "$1 , $2"
    curl -X DELETE -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d "{\"droplet_ids\": [$1]}" "https://api.digitalocean.com/v2/load_balancers/$balancer/droplets"
}

destroy_droplet() 
{
    curl -X DELETE -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" "https://api.digitalocean.com/v2/droplets/$1"
}

get_balancer

echo "Which droplet would you like to remove from the balancer?"
read droplet
delete_droplet_from_balancer $droplet

echo "Would you like to destroy the droplet? (Y/n)"
read userResp
if [ $userResp == "Y" ]
    then
        destroy_droplet $droplet
fi