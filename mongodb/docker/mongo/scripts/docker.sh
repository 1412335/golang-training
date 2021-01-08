#!/bin/bash

# NOTE: This is the simplest way of achieving a replicaset in mongodb with Docker.
# However if you would like a more automated approach, please see the setup.sh file and the docker-compose file which includes this startup script.

# run this after setting up the docker-compose This will instantiate the replica set.
# The id and hostname's can be tailored to your liking, however they MUST match the docker-compose file above.
docker-compose up -d
docker-compose exec mongo1 mongo --port 9142

rs.reconfig(
  {
    _id : 'rs0',
    members: [
      { _id : 0, host : "mongo1:9142" },
      { _id : 1, host : "mongo2:9242" },
      { _id : 2, host : "mongo3:9342", arbiterOnly: true }
    ]
  }
)

exit