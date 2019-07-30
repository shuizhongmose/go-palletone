#!/bin/bash


# now is in bdd/digital_identity directory
pkill gptn

# change to go-palletone directory
cd ../../
cd ./common/crypto

# replace crypto algorithm
sed -i 's/CryptoS256/CryptoP256/g' crypto.go

# edit digital-identity package config file caconfig.yaml
cd $GOPATH/src/github.com/palletone/digital-identity/config
sed -i 's/^url:*$/http:\/\/localhost:7064/g' caconfig.yaml 

# compile gptn
cd $GOPATH/src/github.com/palletone/go-palletone
make gptn
rm -rf bdd/digital_identity/node
mkdir -p bdd/digital_identity/node
cp build/bin/gptn bdd/digital_identity/node

# new genesis
cd bdd/digital_identity/node
chmod +x ./gptn
./gptn newgenesis "" fasle << EOF
y
1
1
EOF

# replace ca certificate bytes
cabytes=`sed ':a;N;$!ba;s/\n/\\n/g' ~/cawork/root/ca-cert.pem`
sed -i "s/^  \"rootCABytes\":*$/  \"rootCABytes\": \"$cabytes\"/g"  ptn-genesis.json


# gptn init
./gptn init << EOF
1
EOF

# start gptn
nohup ./gptn &

