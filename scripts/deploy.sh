#!/bin/bash

# docker on 16.04 id: 23755729
# POST request to set up a new droplet
create_api_server()
{
    curl -X POST 
    -H "Content-Type: application/json" 
    -H "Authorization: Bearer $TOKEN" 
    -d '{
        "name":"$1",
        "region":"sfo2",
        "size":"512mb",
        "image":"docker-16-04",
        "ssh_keys":[7926321],
        "backups":false,
        "ipv6":true,
        "user_data":null,
        "private_networking":null,
        "volumes": null,
        "tags":["api"]
        }' "https://api.digitalocean.com/v2/droplets"
}

create_web_server()
{
    curl `-X POST 
    -H "Content-Type: application/json" 
    -H "Authorization: Bearer $TOKEN" 
    -d '{
        "name":"$1",
        "region":"sfo2",
        "size":"512mb",
        "image":"nginx",
        "ssh_keys":[7926321],
        "backups":false,
        "ipv6":true,
        "user_data":null,
        "private_networking":null,
        "volumes": null,
        "tags":["web"]
        }' "https://api.digitalocean.com/v2/droplets"`
}

echo "Enter the type of server you would like to create (api/web)"
read typ
echo "Enter the name of your server"
read name
#echo "How many servers would you like to create?"
#read num

if [ $typ == "api" ]; then
    create_api_server $name
else if [ $typ == "web" ]; then
    create_web_server $name
    fi
fi

