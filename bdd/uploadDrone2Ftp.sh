#!/bin/bash
#!/usr/bin/expect

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
echo $folder
echo $number
set timeout 120
#set ftppwd [lindex $argv 1]
ftppwd=$1
lftp travis:$ftppwd@47.74.209.46 << EOF
cd $folder
mkdir $number
cd $number
mirror -R /drone/src/github.com/palletone/go-palletone/bdd/logs
exit
EOF