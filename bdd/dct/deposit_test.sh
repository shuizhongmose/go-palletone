#!/bin/bash

./test_case_setup.sh $1

sleep 1

listAccounts=`../node/gptn --exec 'personal.listAccounts' attach ../node/palletone/gptn.ipc` 
key=`echo $listAccounts `
echo $key

list=`echo $key | jq ''`;

length=`echo $key | jq 'length'`

num=$[$length - 1]

for index in `seq 0 $num`
do
	if [ $index -gt 5 ]
	then
	account0=`echo $list|jq ".[0]"`
	another=`echo $list|jq ".[$index]"`
	./transfertoken.sh $account0 $another
	sleep 3 
	fi
#	#echo $list | jq ".[$index]";
done

mediatorAddr_01=`echo $list|jq ".[0]"`

#foundationAddr=`echo $list|jq ".[6]"`

../node/gptn --exec "personal.unlockAccount($mediatorAddr_01,'1',0)" attach ../node/palletone/gptn.ipc

for index in `seq 0 $num`
do
	if [ $index -gt 5 ]
	then
	toAddr=`echo $list | jq ".[$index]"`
	../node/gptn --exec "personal.unlockAccount($toAddr,'1',0)" attach ../node/palletone/gptn.ipc
	fi
	#echo $list | jq ".[$index]";
done

#mdi_01=`echo ${mediatorAddr_01//\"/}`

found=`echo ${mediatorAddr_01//\"/}`

mediatorAddr_02=`echo $list | jq ".[6]"`
mdi_02=`echo ${mediatorAddr_02//\"/}`

#mdiatorAddr_03=`echo $list | jq ".[4]"`
#mdi_03=`echo ${mediatorAddr_03//\"/}`

juryAddr_01=`echo $list | jq ".[7]"`
jury_01=`echo ${juryAddr_01//\"/}`

#获取陪审团节点的公钥
jury_01_pubkey=`../node/gptn --exec "personal.getPublicKey($juryAddr_01)" attach ../node/palletone/gptn.ipc`
jury_01_pub=`echo ${jury_01_pubkey//\"/}`
echo "jury_01_pubkey =>"  $jury_01_pub

developerAddr_01=`echo $list | jq ".[8]"`
developer_01=`echo ${developerAddr_01//\"/}`

juryA_02=`echo $list | jq ".[9]"`
jury_02=`echo ${juryA_02//\"/}`
#获取陪审团节点的公钥
jury_02_pubkey=`../node/gptn --exec "personal.getPublicKey($juryA_02)" attach ../node/palletone/gptn.ipc`
jury_02_pub=`echo ${jury_02_pubkey//\"/}`
echo "jury_02_pubkey =>"  $jury_02_pub

deveA_02=`echo $list | jq ".[10]"`
developer_02=`echo ${deveA_02//\"/}`
m01=`echo $list | jq ".[11]"`
mdi_01=`echo ${m01//\"/}`
m1=`../node/gptn --exec "personal.getPublicKey($m01)" attach ../node/palletone/gptn.ipc`
m11=`echo ${m1//\"/}`
echo "mdi_01_pubkey =>" $m11

m2=`../node/gptn --exec "personal.getPublicKey($mediatorAddr_02)" attach ../node/palletone/gptn.ipc`
m22=`echo ${m2//\"/}`

votedAddr=`echo $list | jq ".[12]"`
votedAddress=`echo ${votedAddr//\"/}`

echo "mdi_02_pubkey =>" $m22
echo "mdi_01" $mdi_01
echo "found" $found
echo "juryA_01" $jury_01
echo "devA_01" $developer_01
echo "mdi_02" $mdi_02
echo "juryA_02" $jury_02
echo "developerAddr_02" $developer_02
echo "votedAddress" $votedAddress
#pybot -d ./log -v mediatorAddr_01:$mdi_01 -v foundationAddr:$found --test Business_01 ./deposit_test_cases/DepositContractTest.robot
#pybot -d ./log -v mediatorAddr_02:$mdi_02 -v foundationAddr:$found --test Business_02 ./deposit_test_cases/DepositContractTest.robot
#pybot -d ./log -v juryAddr_01:$jury_01 -v foundationAddr:$found --test Business_03 ./deposit_test_cases/DepositContractTest.robot
#pybot -d ./log -v developerAddr_01:$developer_01 -v foundationAddr:$found --test Business_04 ./deposit_test_cases/DepositContractTest.robot
#pybot -d ./log -v juryAddr_02:$jury_02 -v foundationAddr:$found --test Business_05 ./deposit_test_cases/DepositContractTest.robot
#pybot -d ./log -v developerAddr_02:$developer_02 -v foundationAddr:$found --test Business_06 ./deposit_test_cases/DepositContractTest.robot
pybot -d ../logs/deposit -v votedAddress:$votedAddress -v m1_pubkey:$m11 -v m2_pubkey:$m22 -v juryAddr_01_pubkey:$jury_01_pub -v juryAddr_02_pubkey:$jury_02_pub -v mediatorAddr_01:$mdi_01 -v foundationAddr:$found -v mediatorAddr_02:$mdi_02 -v juryAddr_01:$jury_01 -v developerAddr_01:$developer_01 -v juryAddr_02:$jury_02 -v developerAddr_02:$developer_02 ./deposit_test_cases/DepositContractTest.robot
#pybot -d ./log -v mediatorAddr_01:$mdi_01 -v foundationAddr:$found -v mediatorAddr_02:$mdi_02 -v juryAddr_01:$jury_01 -v developerAddr_01:$developer_01 -v juryAddr_02:$jury_02 -v developerAddr_02:$developer_02 --test Business_01 --test Business_03 --test Business_05 ./deposit_test_cases/DepositContractTest.robot
#pybot -d ./log -v mediatorAddr_01:$mdi_01 -v foundationAddr:$found -v otherAddr:$otherAddr --test Business_07 ./deposit_test_cases/DepositContractTest.robot

#pybot -d ./log -v foundationAddr:$found  -v votedAddress:$votedAddress --test PledgeTest ./deposit_test_cases/DepositContractTest.robot

#./test_case_teardown.sh

#killall gptn
