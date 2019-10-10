#!/bin/bash

# not exit while error
set +e

# project environment
# you can use BASE_DIR, it points $GOPATH/src/github.com/palletone/go-palletone
export ALL_LOG_PATH=$GOPATH/src/github.com/palletone/go-palletone/bdd/node/log/all.log
export GAS_TOKEN_ALL_LOG_PATH=$GOPATH/src/github.com/palletone/go-palletone/bdd/GasToken/node/log/all.log
export BDD_LOG_PATH=$GOPATH/src/github.com/palletone/go-palletone/bdd/logs
export CREATE_TRANS_DIR=createTrans
export CONTRACT20_DIR=crt20Contract
export SEQENCE721_DIR=crt721Seqence
export UDID721_DIR=crt721UDID
export VOTECONTRACT_DIR=voteContract
export MULTIPLE_DIR=zMulti-node
export DIGITAL_IDENTITY_DIR=Digital-Identitycert
export DEPOSIT_DIR=deposit
export GAS_TOKEN_DIR=gasToken
export MEDIATOR_VOTE_DIR=meidatorvote
export USER_CONTRACT_DIR=usercontract
export GO111MODULE=on
export FTP_PWD=Pallet2018

# before install
apt-get install python2.7 -y
wget https://bootstrap.pypa.io/get-pip.py
python get-pip.py
pip -V
pip install --upgrade pip
pip install robotframework==2.8.5
pip install requests
pip install robotframework-requests
pip install demjson
#pip install pexpect
apt-get install tcl tk expect -y
whereis expect
apt-get install jq -y
apt-get install lftp -y
apt-get install tree -y
apt-get install net-tools -y
cd $BASE_DIR
chmod +x bdd/uploadDrone2Ftp.sh
chmod +x bdd/upload2Ftp.sh

# run tests
cd $BASE_DIR/bdd/${test_dir}/scripts
chmod +x start.sh
./start.sh

# upload logs
#chmod +x upload.sh
#./upload.sh

echo "upload files done"

# after tests
pkill gptn
pkill fabric