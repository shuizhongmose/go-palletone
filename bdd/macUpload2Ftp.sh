#!/usr/bin/expect
#!/bin/bash
set timeout 120
set ftppwd [lindex $argv 0]
set folder [lindex $argv 1]
spawn ftp 47.74.209.46
expect "Name"
send "travis\r"
expect "Password"
send "${ftppwd}\r"
expect "successful"
send "cd ${folder}\r"
expect "changed"
send "put ./gptn-mac.tar.gz ./go-palletone_mainnet_v1.0.0_macOS-amd64.tar.gz\r"  
expect "Ok"
send "bye\r"
interact