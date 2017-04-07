#!/bin/bash

# docker on 16.04 id: 23755729
# POST request to set up a new droplet
create_server()
{
    body=$(cat  << EOF
    {
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
    }
EOF
)
    curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $TOKEN" -d "$body"  "https://api.digitalocean.com/v2/droplets"
}


echo "Enter the name of your server"
read name

echo "How many servers would you like to create?"
read num

# if [ $typ == "api" ]; then
#     create_api_server $name
# else if [ $typ == "web" ]; then
#     create_web_server $name
#     fi
# fi

