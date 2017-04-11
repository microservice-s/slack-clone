#!/bin/bash

# docker on 16.04 id: 23755729

# POST request to set up a new droplet
create_server()
{
    user_data=$(cat ./user-data.yml)
    body=$(cat  <<EOF
    {
        "name":"$1",
        "region":"sfo2",
        "size":"512mb",
        "image":"docker-16-04",
        "user_data":$user_data
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
    # create a new droplet with the given name and token env var
    # then get its id using jq
    resp=$(curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $TOKEN" -d "$body"  "https://api.digitalocean.com/v2/droplets")
    id=$(echo $resp | jq -r .droplet.id)
    # add it to the load balancer
    status=$(echo $resp | jq -r .droplet.status)

    # keep checking the status and add it to the load balancer when it's ready
    while [ "$status" != "active" ]; do
        # sleep for 10 seconds and try again
        sleep 10
        echo "Checking status of $id"
        status=$(curl -X GET -H "Content-Type: application/json" -H "Authorization: Bearer $TOKEN" "https://api.digitalocean.com/v2/droplets/$id" | jq -r .droplet.status)
        echo "Status is: $status"
    done

    add_to_balance $id
}

test_server()
{
    user_data=$(cat ./user-data.yml)
    echo $user_data
    body=$(cat  <<EOF
    {
        "name":"$1",
        "region":"sfo2",
        "size":"512mb",
        "image":"docker-16-04",
        "user_data":"$user_data",
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
    # create a new droplet with the given name and token env var
    # then get its id using jq
    curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $TOKEN" -d "$body"  "https://api.digitalocean.com/v2/droplets"
}

# take a given droplet id and put it behind a load balancer as specified in the env var $LOAD_BALANCER_ID
add_to_balance()
{
    body=$(cat <<EOF
    {
        "droplet_ids": [$1]
    }
EOF
)
    curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $TOKEN" -d "$body" "https://api.digitalocean.com/v2/load_balancers/$LOAD_BALANCER_ID/droplets"
}


echo "Enter the name of your server"
read name
test_server $name
# # echo "How many servers would you like to create?"
# # read num
#create_server $name



#add_to_balance 45423708

# if [ $typ == "api" ]; then
#     create_api_server $name
# else if [ $typ == "web" ]; then
#     create_web_server $name
#     fi
# fi

