#!/bin/bash

function obtain_git_branch {
  br=`git branch | grep "*"`
  echo ${br/* /}
}

function obtain_git_commit_number {
  br=`git rev-list --count HEAD`
  echo ${br/ /}
}

folder=`obtain_git_branch`
number=`obtain_git_commit_number`
ftppwd=$1

LOG_PATH=$BDD_LOG_PATH
echo $folder
echo $number
lftp travis:$ftppwd@47.74.209.46 << EOF
cd ${folder}
mkdir ${number}
cd ${number}
mirror -R  /home/jenkinsgo/src/github.com/palletone/go-palletone/bdd/logs
exit
EOF