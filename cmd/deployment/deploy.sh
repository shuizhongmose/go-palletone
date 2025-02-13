#!/bin/bash

source ./modifyconfig.sh

function ExecInit()
{
    count=1  
    while [ $count -le $1 ] ;  
    do  
    #echo $count;
    #ExecDeploy $count 
    if [ $count -eq 1 ] ;then
    cd node$count
    cp ../init.sh .
    gptninit=`./init.sh`
    initinfo=`echo $gptninit | sed -n '$p'`
    initinfotemp=`echo $initinfo | awk '{print $NF}'`
    initinfotemp=${initinfotemp:0:7}
       if [ $initinfotemp != "success" ] ; then
               echo "====================init err=================="
               echo $initinfo
               return
       fi

    path=`pwd`
    fullpath=${path}"/palletone/leveldb"
    #echo "leveldb path:"$fullpath
    if [ ! -d $fullpath ]; then
        echo "====================init err=================="
        return
    fi
        rm -rf init.sh log
        cd ../
    else
        cd node$count
        cp ../node1/palletone/leveldb ./palletone/. -rf
        rm -rf log
        cd ../
    fi
    let ++count;  
    sleep 1;  
    done

    length=${#initinfo}
    num=$[$length-112]
    str=${initinfo:$num:112}
    charToSearch="\[";
    let pos=`echo "$str" | awk -F ''$charToSearch'' '{printf "%d", length($0)-length($NF)}'`
    genesishash=${str:$pos:66}
    echo "Init OK GenesisHash="$genesishash

    return 0;  
}



function replacejson()
{
    length=`cat $1 |jq '.initialMediatorCandidates| length'`

    add=`cat $1 | jq ".initialParameters.active_mediator_count = $length"`

    add=`echo $add | jq ".immutableChainParameters.min_mediator_count = $length"`

    add=`echo $add | jq ".initialParameters.maintenance_skip_slots = 2"`

    add=`echo $add | jq ".immutableChainParameters.min_maint_skip_slots = 2"`

    add=`echo $add | jq ".initialParameters.mediator_interval = 3"`

    if [ -n "$2" ]; then
        add=`echo $add | jq ".initialParameters.contract_election_num = 3"`
    fi

    tempstamp=`cat $1 | jq '.initialTimestamp'`
    tempstamp=$[$tempstamp/3]
    tempstamp=`echo $tempstamp | cut -f1 -d"."`
    tempstamp=$[$tempstamp*3]

    add=`echo $add | jq ".initialTimestamp = $tempstamp"`

    rm $1
    echo $add >> temp.json
    jq -r . temp.json >> $1
    rm temp.json
    echo "======modify genesis json ok======="
}

function ExecDeploy()
{
    echo ===============$1===============
    mkdir "node"$1
    #if [ $1 -eq 4 ] ;then
    #echo "=="$1
    cp gptn ./createaccount.sh modifyconfig.sh modifyjson.sh init.sh createpk.sh node$1
    cd node$1
    /bin/bash modifyconfig.sh
    #source ./modifyjson.sh
    source ./modifyconfig.sh
    ModifyConfig $1
    rm *.sh
    cd ../
    #else
    #echo $1
    #fi
}



function LoopDeploy()  
{  
    count=1;  
    while [ $count -le $1 ] ;  
    do  
    #echo $count;
    ExecDeploy $count 
    let ++count;  
    sleep 1;  
    done  
    return 0;  
}
path=`echo $GOPATH`
src=/src/github.com/palletone/go-palletone/build/bin/gptn
fullpath=$path$src
cp $fullpath .

n=
if [ -n "$1" ]; then
    n=$1
else
    read -p "Please input the numbers of nodes you want: " n;
fi

eleNum=
if [ -n "$2" ]; then
    eleNum=$2
fi

LoopDeploy $n;

json="node1/ptn-genesis.json"
replacejson $json $eleNum

#ModifyBootstrapNodes $n

#ExecInit $n
initvalue=$(ExecInit $n)
echo $initvalue
charToSearch="GenesisHash=";
let pos=`echo "$initvalue" | awk -F ''$charToSearch'' '{printf "%d", length($0)-length($NF)}'`
genesishash=${initvalue:$pos:66}

ModifyP2PConfig $n $genesishash

num=$[$n+1]
MakeTestNet $num $genesishash

num=$[$n+2]
MakeTestNet $num $genesishash

num=$[$n+3]
MakeTestNet $num $genesishash


num=$[$n+4]
MakeTestNet $num $genesishash

num=$[$n+5]
MakeTestNet $num $genesishash

rm -f node*/gptn
