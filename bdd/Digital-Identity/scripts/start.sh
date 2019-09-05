#!/bin/bash

cd $GOPATH/src/github.com/
cd palletone
git clone https://github.com/palletone/digital-identity.git
cd $BASE_DIR/bdd/Digital-Identity/scripts
chmod +x ca-start.sh
chmod +x one-ptn.sh
./ca-start.sh
tree ~/cawork
./one-ptn.sh
sleep 120
netstat -ntl
cd $BASE_DIR
mkdir -p ${BDD_LOG_PATH}/${DIGITAL_IDENTITY_DIR}
python -m robot.run -d ${BDD_LOG_PATH}/${DIGITAL_IDENTITY_DIR} ./bdd/Digital-Identity/testcases
