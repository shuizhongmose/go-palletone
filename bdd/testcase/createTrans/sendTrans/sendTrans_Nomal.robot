*** Settings ***
Default Tags      normal
Library           RequestsLibrary
Library           Collections
Library           ../../utilFunc/createToken.py
Resource          ../../utilKwd/utilDefined.txt
Resource          ../../utilKwd/behaveKwd.txt
Resource          ../../utilKwd/normalKwd.txt

*** Variables ***
${host}           http://localhost:8545/
${geneAdd}        P1CwGYGSjWSaJrysHAjAWtDyFSsbcYwoULv
${recieverAdd}    P11LKXXsDuKUdo3cJEy7uMnqnvw4kwHPs8
${givenAmount}    10
@{transList}      ${geneAdd}    ${recieverAdd}    ${givenAmount}    2
${PTN}            \d+
${result_code}    \f[a-z0-9]*
${result_hex}     \f[a-z0-9]*
${result_txid}    \0[a-z0-9]{60,70}
${sendResult}     [a-z0-9]*

*** Test Cases ***
sendTransNormal
    [Tags]    normal
    ${PTN1}    ${result1}    normalGetBalance    ${recieverAdd}
    ${result11}    Evaluate    ${PTN1}+${givenAmount}
    ${sendResult}    normalSendTrans
    Sleep    4
    ${PTN2}    ${result2}    normalGetBalance    ${recieverAdd}
    Sleep    1
    Should Be Equal As Strings    ${result11}    ${PTN2}
