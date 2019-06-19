*** Settings ***
Default Tags      nomal
Library           ../../utilFunc/createToken.py
Resource          ../../utilKwd/utilVariables.txt
Resource          ../../utilKwd/normalKwd.txt
Resource          ../../utilKwd/utilDefined.txt
Resource          ../../utilKwd/behaveKwd.txt

*** Variables ***
${preTokenId}     QA086

*** Test Cases ***
Scenario: 721 Contract - Supply token
    [Documentation]    Verify Sender's PTN and token
    #${ret}    Given CcinvokePass normal
    Given CcinvokePass normal
    ${PTN1}    And Request getbalance before create token
    ${ret}    When Spply token of 721 contract
    ${PTNGAIN}    Calculate gain
    ${PTN2}    Request getbalance after transfer token
    Then Assert gain    ${PTN1}    ${PTN2}    ${PTNGAIN}

*** Keywords ***
CcinvokePass normal
    ${geneAdd}    getGeneAdd    ${host}
    Set Suite Variable    ${geneAdd}    ${geneAdd}
    ${ccList}    Create List    ${crtTokenMethod}    ${note}    ${preTokenId}    ${UDIDToken}    ${721TokenAmount}
    ...    ${721MetaBefore}    ${geneAdd}
    ${resp}    Request CcinvokePass    ${commonResultCode}    ${geneAdd}    ${recieverAdd}    ${PTNAmount}    ${PTNPoundage}
    ...    ${721ContractId}    ${ccList}
    sleep    4
    ${jsonRes}    Evaluate    demjson.encode(${resp.content})    demjson
    ${jsonRes}    To Json    ${jsonRes}
    [Return]    ${jsonRes['result']}

Request getbalance before create token
    #${PTN1}    ${result1}    normalGetBalance    ${geneAdd}
    ${result1}    getBalance    ${geneAdd}
    sleep    4
    ${PTN1}    Get From Dictionary    ${result1}    PTN
    #sleep    1
    #${coinToken1}    Get From Dictionary    ${result1}    ${key}
    [Return]    ${PTN1}

Spply token of 721 contract
    ${ccList}    Create List    ${supplyTokenMethod}    ${preTokenId}    ${721TokenAmount}    ${721MetaAfter}
    ${resp}    Request CcinvokePass    ${commonResultCode}    ${geneAdd}    ${recieverAdd}    ${PTNAmount}    ${PTNPoundage}
    ...    ${721ContractId}    ${ccList}
    ${jsonRes}    Evaluate    demjson.encode(${resp.content})    demjson
    ${jsonRes}    To Json    ${jsonRes}
    [Return]    ${jsonRes['result']}

Calculate gain
    ${PTNGAIN}    Evaluate    ${PTNAmount}+${PTNPoundage}
    ${PTNGAIN}    countRecieverPTN    ${PTNGAIN}
    sleep    2
    [Return]    ${PTNGAIN}

Request getbalance after transfer token
    #normalCcqueryById    ${721ContractId}    getTokenInfo    ${preTokenId}
    ${PTN2}    ${result2}    normalGetBalance    ${geneAdd}
    sleep    4
    #${key}    getTokenId    ${preTokenId}    ${result2['result']}
    #${queryResult}    ccqueryById    ${721ContractId}    ${existToken}    ${key}
    #Should Be Equal As Strings    ${queryResult['result']}    True
    ${queryResult}    ccqueryById    ${721ContractId}    ${TokenInfoMethod}    ${preTokenId}
    ${tokenCommonId}    ${countList}    jsonLoads    ${queryResult['result']}    AssetID    TokenIDs
    ${len}    Evaluate    len(${countList})
    Should Be Equal As Numbers    ${len}    10
    : FOR    ${num}    IN RANGE    ${len}
    \    ${number}    Evaluate    ${num}+1
    \    ${key}    getTokenIdByNum    ${tokenCommonId}    ${result2['result']}    ${number}
    \    ${voteToken}    Get From Dictionary    ${result2['result']}    ${key}
    \    log    ${key}
    \    Should Be Equal As Numbers    ${voteToken}    1
    [Return]    ${PTN2}

Assert gain
    [Arguments]    ${PTN1}    ${PTN2}    ${PTNGAIN}
    #${result2}    getBalance    ${geneAdd}
    #sleep    4
    #${PTN2}    Get From Dictionary    ${result2}    PTN
    ${GAIN}    Evaluate    decimal.Decimal('${PTN1}')-decimal.Decimal('${PTNGAIN}')    decimal
    Should Be Equal As Numbers    ${PTN2}    ${GAIN}
