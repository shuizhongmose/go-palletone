#!/bin/bash


pkill fabric

# install ca
export GO111MODULE=off
cd $GOPATH/src
go get -u github.com/hyperledger/fabric-ca/cmd/...
cd $GOPATH/src/github.com/hyperledger/fabric-ca/
make fabric-ca-server
export PATH=$GOPATH/src/github.com/hyperledger/fabric-ca/bin:$PATH

cd ~
rm -rf cawork
mkdir cawork
cd cawork
mkdir root immediateca

# 初始化根CA
cd root
fabric-ca-server init -b lk:123

# 修改配置文件 fabric-ca-server-config.yaml
num=`grep  -n  "org2:"   fabric-ca-server-config.yaml  | head -1  | cut  -d  ":"  -f  1`
let nextnum=$num+1
sed -i "${num},${nextnum}d" fabric-ca-server-config.yaml

sed -i 's/org1:/gptn:/g'  fabric-ca-server-config.yaml
sed -i 's/- department1/- mediator1/g'  fabric-ca-server-config.yaml
sed -i 's/- department2/- mediator2/g'  fabric-ca-server-config.yaml

# 启动CA
nohup fabric-ca-server start -b lk:123 >> caserver.out &

sleep 10
# 进入immediateca
cd ../immediateca
# initialize immediate ca
fabric-ca-server init -b lk:123 -u http://lk:123@localhost:7054
# edit config file
num=`grep  -n  "org2:"   fabric-ca-server-config.yaml  | head -1  | cut  -d  ":"  -f  1`
let nextnum=$num+1
sed -i "${num},${nextnum}d" fabric-ca-server-config.yaml

sed -i 's/org1:/gptn:/g'  fabric-ca-server-config.yaml
sed -i 's/- department1/- mediator1/g'  fabric-ca-server-config.yaml
sed -i 's/- department2/- mediator2/g'  fabric-ca-server-config.yaml
# change server port
sed -i 's/9443/9453/g' fabric-ca-server-config.yaml
# start immediate ca
nohup fabric-ca-server start -b lk:123 -p 7064 -u http://lk:123@localhost:7054 >> immediate.out &

