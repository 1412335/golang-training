#!/bin/bash

#MONGODB1=`ping -c 1 mongo1 | head -1  | cut -d "(" -f 2 | cut -d ")" -f 1`
#MONGODB2=`ping -c 1 mongo2 | head -1  | cut -d "(" -f 2 | cut -d ")" -f 1`
#MONGODB3=`ping -c 1 mongo3 | head -1  | cut -d "(" -f 2 | cut -d ")" -f 1`

MONGODB1=mongo1:9142
MONGODB2=mongo2:9242
MONGODB3=mongo3:9342

echo "**********************************************" ${MONGODB1}
echo "Waiting for startup.."
until curl http://${MONGODB1}/serverStatus\?text\=1 2>&1 | grep uptime | head -1; do
  printf '.'
  sleep 1
done

# echo curl http://${MONGODB1}/serverStatus\?text\=1 2>&1 | grep uptime | head -1
# echo "Started.."


echo SETUP.sh time now: `date +"%T" `
# mongo --username root --password root --authenticationDatabase admin --host ${MONGODB1} <<EOF
mongo --host ${MONGODB1} <<EOF
var cfg = {
    "_id": "rs0",
    "protocolVersion": 1,
    "version": 1,
    "members": [
        {
            "_id": 0,
            "host": "${MONGODB1}",
            "priority": 2
        },
        {
            "_id": 1,
            "host": "${MONGODB2}",
            "priority": 0
        },
        {
            "_id": 2,
            "host": "${MONGODB3}",
            "priority": 0
        }
    ],settings: {chainingAllowed: true}
};
rs.initiate(cfg, { force: true });
rs.reconfig(cfg, { force: true });
rs.secondaryOk();
db.getMongo().setReadPref('nearest');
db.getMongo().setSecondaryOk(); 
rs.status();
exit;
EOF

sleep 100;
mongo --host ${MONGODB1} <<EOF
show dbs;

use admin;
db.createUser({	
    user: "root",
	pwd: "root",
	roles:[{role: "userAdminAnyDatabase" , db:"admin"}]
});

use go_mongo;
db.createCollection("article_category");
show collections;
db.createUser({
	user: "root",
	pwd: "12345",
	roles:[{role: "userAdmin" , db:"go_mongo"}]
});
show users;
EOF