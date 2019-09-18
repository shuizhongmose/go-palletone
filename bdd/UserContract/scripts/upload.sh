#!/bin/bash

cd ${BASE_DIR}
zip -j ./bdd/logs/uc3-3_reports.zip ${BDD_LOG_PATH}/${USER_CONTRACT_DIR}/*
rm -rf ${BDD_LOG_PATH}/${USER_CONTRACT_DIR}
zip -r -l -o  ./bdd/logs/uc3-3_logs.zip ./cmd/deployment/node1/log ./cmd/deployment/node1/nohup.out ./cmd/deployment/node2/log ./cmd/deployment/node2/nohup.out ./cmd/deployment/node3/log ./cmd/deployment/node3/nohup.out
ls ./bdd/logs
./bdd/uploadDrone2Ftp.sh ${FTP_PWD}