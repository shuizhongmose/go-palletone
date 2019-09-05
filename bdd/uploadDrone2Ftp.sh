#!/usr/bin/expect
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

set timeout 120
set ftppwd [lindex $argv 0]
spawn lftp travis:$ftppwd@47.74.209.46
expect "lftp"
send "cd ${folder}\n"
expect "cd"
send "mkdir ${number}\n"
expect "mkdir"
send "cd ${number}\n"
expect "cd"
send "mirror -R /home/travis/gopath/src/github.com/palletone/go-palletone/bdd/logs\n"
expect "transferred"
send "exit\n"
interact