#!/bin/bash

cd $BASE_DIR
make gptn
make golang-baseimage-dev
docker images
cat /proc/sys/net/ipv4/ip_forward
chmod +x jurySec.sh
./jurySec.sh
cat /etc/docker/daemon.json
cd ./cmd/deployment
rm -rf ./node*
chmod +x ./deploy.sh
chmod +x ./start.sh
cp $BASE_DIR/bdd/UserContract/scripts/editConfig.sh ./
chmod +x editConfig.sh
docker network ls
./deploy.sh 3 3
./editConfig.sh 3 3
./start.sh 3
sleep 30
docker network ls
cd $BASE_DIR
mkdir -p ${BDD_LOG_PATH}/${USER_CONTRACT_DIR}
python -m robot.run -d ${BDD_LOG_PATH}/${USER_CONTRACT_DIR} ./bdd/UserContract/Testcases;
docker ps -a