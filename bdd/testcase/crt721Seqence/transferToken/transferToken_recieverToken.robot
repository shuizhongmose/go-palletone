*** Settings ***
Default Tags      normal
Library           ../../utilFunc/createToken.py
Resource          ../../utilKwd/utilVariables.txt
Resource          ../../utilKwd/normalKwd.txt
Resource          ../../utilKwd/utilDefined.txt
Resource          ../../utilKwd/behaveKwd.txt

*** Variables ***
${preTokenId}     QA077
${subId}          4

*** Test Cases ***
Scenario: 721 Contract - Transfer token
    [Documentation]    Verify Reciever's Token
    Given Get genesis address
    ${ret}    When Create token of 721 contract
    ${key}    ${voteToken}    And Request getbalance before transfer token
    And Request transfer token    ${key}
    ${voteToken2}    And Request getbalance after transfer token    ${key}
    Then Assert gain    ${voteToken}    ${voteToken2}

*** Keywords ***
Get genesis address
    ${geneAdd}    getGeneAdd    ${host}
    Set Suite Variable    ${geneAdd}    ${geneAdd}
    personalUnlockAccount    ${geneAdd}
<<<<<<< HEAD
    sleep    4
=======
>>>>>>> master

Create token of 721 contract
    ${ccList}    Create List    ${crtTokenMethod}    ${note}    ${preTokenId}    ${SeqenceToken}    ${721TokenAmount}
    ...    ${721MetaBefore}    ${geneAdd}
    ${resp}    Request CcinvokePass    ${commonResultCode}    ${geneAdd}    ${recieverAdd}    ${PTNAmount}    ${PTNPoundage}
    ...    ${721ContractId}    ${ccList}
    ${jsonRes}    Evaluate    demjson.encode(${resp.content})    demjson
    ${jsonRes}    To Json    ${jsonRes}
    [Return]    ${jsonRes['result']}

Request getbalance before transfer token
    sleep    4
    ${PTN1}    ${result1}    normalGetBalance    ${geneAdd}
<<<<<<< HEAD
    sleep    2
=======
>>>>>>> master
    ${queryResult}    ccqueryById    ${721ContractId}    getTokenInfo    ${preTokenId}
    ${tokenCommonId}    ${countList}    jsonLoads    ${queryResult['result']}    AssetID    TokenIDs
    Set Suite Variable    ${key}    ${tokenCommonId}-${subId}
    ${voteToken}    Get From Dictionary    ${result1['result']}    ${key}
    [Return]    ${key}    ${voteToken}

Request transfer token
    [Arguments]    ${key}
    ${tokenResult}    transferToken    ${key}    ${geneAdd}    ${recieverAdd}    1    ${PTNPoundage}
    ...    ${evidence}    ${duration}


Request getbalance after transfer token
    [Arguments]    ${key}
    sleep    4
    ${PTN1}    ${result2}    normalGetBalance    ${recieverAdd}
<<<<<<< HEAD
    sleep    2
=======
>>>>>>> master
    ${voteToken2}    Get From Dictionary    ${result2['result']}    ${key}
    [Return]    ${voteToken2}

Assert gain
    [Arguments]    ${voteToken}    ${voteToken2}
    Should Be Equal As Strings    ${voteToken}    ${voteToken2}
