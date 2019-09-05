#!/bin/bash

cd /drone/src/github.com/palletone/go-palletone
cp bdd/Digital-Identity/node/log/all.log $BDD_LOG_PATH/$DIGITAL_IDENTITY_DIR
cp bdd/Digital-Identity/node/nohup.out $BDD_LOG_PATH/$DIGITAL_IDENTITY_DIR
ls $BDD_LOG_PATH/$DIGITAL_IDENTITY_DIR
echo ${FTP_PWD}
./bdd/uploadDrone2Ftp.sh ${FTP_PWD}
