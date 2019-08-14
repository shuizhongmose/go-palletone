#!/bin/bash


# now is in bdd/Digital-Identity directory
pkill gptn

# change to go-palletone directory
cd $GOPATH/src/github.com/palletone/go-palletone
git pull
cd ./common/crypto

# replace crypto algorithm
sed -i 's/CryptoS256/CryptoP256/g' crypto.go

# edit digital-identity package config file caconfig.yaml
cd $GOPATH/src/github.com/palletone/digital-identity/config
sed -i 's/^url:.*$/url: http:\/\/localhost:7064/g' caconfig.yaml 

# compile gptn
cd $GOPATH/src/github.com/palletone/go-palletone
export GO111MODULE=on
make gptn
rm -rf bdd/Digital-Identity/node
mkdir -p bdd/Digital-Identity/node
cp build/bin/gptn bdd/Digital-Identity/node

# new genesis
cd bdd/Digital-Identity/node
chmod +x ./gptn
./gptn newgenesis "" fasle << EOF
y
1
1
EOF

# replace ca certificate bytes
cabytes=`sed ':a;N;$!ba;s/\n/\\n/g' ~/cawork/root/ca-cert.pem`
echo $cabytes
res=`cat ptn-genesis.json | jq ".digitalIdentityConfig.rootCABytes= \"$cabytes\""`
echo $res >> tmp.json
mv ptn-genesis.json ptn-genesis.org.json
jq -r . tmp.json >> ptn-genesis.json
rm tmp.json


#change http port
sed -i "s/8545/8645/g" ptn-config.toml
sed -i "s/HTTPHost = \"localhost\"/HTTPHost = \"0.0.0.0\"/g" ptn-config.toml

# gptn init
./gptn init << EOF
1
EOF

# start gptn
nohup ./gptn &

