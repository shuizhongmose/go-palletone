##!/usr/bin/expect
#!/bin/bash
#set timeout 30
#spawn ../node/gptn account new
#expect "Passphrase:"
#send "1\n"
#expect "Repeat passphrase:"
#send "1\n"
#interact
./gptn account new << EOF
1
1
EOF
