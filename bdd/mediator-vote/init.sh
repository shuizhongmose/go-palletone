#!/bin/bash
#pkill gptn
#tskill gptn
#cd ../../cmd/gptn && go build
cd ../../
#rm -rf ./bdd/GasToken/node/*
#cp ./cmd/gptn/gptn ./bdd/GasToken/node
cd ./bdd/GasToken/node
chmod +x gptn

# new genesis
./gptn newgenesis "" fasle << EOF
y
1
1
EOF

# edit genesis json
gasToken="WWW"
jsonFile="ptn-genesis.json"
if [ -e "$jsonFile" ]; then
    #file already exist, modify
    sed -i "s/\"gasToken\": \"PTN\"/\"gasToken\": \"$gasToken\"/g" $jsonFile
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