#!/bin/bash
#pkill gptn
#tskill gptn
#cd ../../cmd/gptn && go build
cd ../../
#rm -rf ./bdd/mediator-vote/node*
#cp ./cmd/gptn/gptn ./bdd/mediator-vote/node
cd ./bdd/mediator-vote/node
chmod +x gptn

# new genesis
./gptn newgenesis "" fasle << EOF
y
1
1
EOF

# edit genesis json
jsonFile="ptn-genesis.json"
if [ -e "$jsonFile" ]; then
    #file already exist, modify
    sed -i "s/\"activeMediatorCount\": \"5\"/\"activeMediatorCount\": \"3\"/g" $jsonFile
    sed -i "s/\"initialActiveMediators\": \"5\"/\"initialActiveMediators\": \"3\"/g" $jsonFile
    sed -i "s/\"minMediatorCount\": \"5\"/\"minMediatorCount\": \"3\"/g" $jsonFile
    sed -i "s/\"maintenanceInterval\": \"600\"/\"maintenanceInterval\": \"150\"/g" $jsonFile
else
    #file not found, new file
    echo "no $jsonFile"
    exit -1
fi

# gptn init
./gptn init << EOF
1
EOF

# start gptn
nohup ./gptn &