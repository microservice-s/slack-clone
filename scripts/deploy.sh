#!/bin/bash

# docker on 16.04 id: 23755729

# POST request to set up a new droplet
create_server()
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
        sleep 15
        echo "Checking status of $id"
        status=$(curl -X GET -H "Content-Type: application/json" -H "Authorization: Bearer $TOKEN" "https://api.digitalocean.com/v2/droplets/$id" | jq -r .droplet.status)
        echo "Status is: $status"
    done

    add_to_balance $id
}

# take a given droplet id and put it behind a load balancer as specified in the env var $LOAD_BALANCER_ID
add_to_balance()
{
    echo "Adding server: $1 to load balancer"
    curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $TOKEN" -d "{\"droplet_ids\": [$1]}" "https://api.digitalocean.com/v2/load_balancers/$LOAD_BALANCER_ID/droplets"
}


echo "Enter the name of your server"
read name

create_server $name

