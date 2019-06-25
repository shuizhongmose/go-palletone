*** Settings ***
Default Tags      normal
Library           ../../utilFunc/createToken.py
Resource          ../../utilKwd/utilVariables.txt
Resource          ../../utilKwd/normalKwd.txt
Resource          ../../utilKwd/utilDefined.txt
Resource          ../../utilKwd/behaveKwd.txt

*** Variables ***
${preTokenId}     QA058

*** Test Cases ***
Scenario: 20Contract- Supply token
    [Documentation]    Verify Sender's PTN
    Given CcinvokePass normal
    ${PTN1}    ${key}    ${coinToken1}    And Request getbalance before create token    ${geneAdd}
    ${ret}    When Create token of vote contract    ${geneAdd}
    ${GAIN}    And Calculate gain of recieverAdd
    ${PTN2}    ${coinToken2}    And Request getbalance after create token    ${geneAdd}    ${key}
    Then Assert gain of reciever    ${PTN1}    ${PTN2}    ${GAIN}    ${coinToken1}    ${coinToken2}

*** Keywords ***
CcinvokePass normal
    ${geneAdd}    getGeneAdd    ${host}
    Set Suite Variable    ${geneAdd}    ${geneAdd}
    ${ccList}    Create List    ${crtTokenMethod}    ${evidence}    ${preTokenId}    ${tokenDecimal}    ${tokenAmount}
    ${ret}    normalCcinvokePass    ${commonResultCode}    ${geneAdd}    ${recieverAdd}    ${PTNAmount}    ${PTNPoundage}
    ...    ${20ContractId}    ${ccList}
<<<<<<< HEAD
    sleep    6
=======
>>>>>>> master

Request getbalance before create token
    [Arguments]    ${geneAdd}
    #${PTN1}    ${result1}    normalGetBalance    ${geneAdd}
    sleep    4
    ${result1}    getBalance    ${geneAdd}
    ${key}    getTokenId    ${preTokenId}    ${result1}
    ${PTN1}    Get From Dictionary    ${result1}    PTN
    ${coinToken1}    Get From Dictionary    ${result1}    ${key}
    sleep    2
    [Return]    ${PTN1}    ${key}    ${coinToken1}

Create token of vote contract
    [Arguments]    ${geneAdd}
    ${ccTokenList}    Create List    ${supplyTokenMethod}    ${preTokenId}    ${supplyTokenAmount}
    ${ccList}    Create List    ${geneAdd}    ${recieverAdd}    ${PTNAmount}    ${PTNPoundage}    ${20ContractId}
    ...    ${ccTokenList}    ${pwd}    ${duration}    ${EMPTY}
    ${resp}    setPostRequest    ${host}    ${invokePsMethod}    ${ccList}
    log    ${resp.content}
    Should Contain    ${resp.content}['jsonrpc']    "2.0"    msg="jsonrpc:failed"
    Should Contain    ${resp.content}['id']    1    msg="id:failed"
    ${ret}    Should Match Regexp    ${resp.content}['result']    ${commonResultCode}    msg="result:does't match Result expression"
    [Return]    ${ret}

Calculate gain of recieverAdd
<<<<<<< HEAD
    sleep    6
    ${invokeGain}    Evaluate    int(${PTNAmount})+int(${PTNPoundage})
    ${GAIN}    countRecieverPTN    ${invokeGain}
    sleep    1
=======
    ${invokeGain}    Evaluate    int(${PTNAmount})+int(${PTNPoundage})
    ${GAIN}    countRecieverPTN    ${invokeGain}
>>>>>>> master
    [Return]    ${GAIN}

Request getbalance after create token
    [Arguments]    ${geneAdd}    ${key}
    sleep    4
    ${result2}    getBalance    ${geneAdd}
    ${coinToken2}    Get From Dictionary    ${result2}    ${key}
    ${PTN2}    Get From Dictionary    ${result2}    PTN
    [Return]    ${PTN2}    ${coinToken2}

Assert gain of reciever
    [Arguments]    ${PTN1}    ${PTN2}    ${GAIN}    ${coinToken1}    ${coinToken2}
    ${PTNGAIN}    Evaluate    decimal.Decimal('${PTN1}')-decimal.Decimal('${GAIN}')    decimal
    Should Be Equal As Numbers    ${PTN2}    ${PTNGAIN}
    Should Be Equal As Numbers    ${coinToken1}    ${coinToken2}
    ${result}    getTxByReqId    ${ret}
    ${jsonRes}    Evaluate    demjson.encode(${result})    demjson
    #${jsonRes}    To Json    ${jsonRes}
    #${TYPE}    Evaluate    str(type(${jsonRes['result']}))
    #log    ${result['item']}
    ${error_code}    Should Match Regexp    ${result}['result']    "error_code":500
    ${error_code}    Should Match Regexp    ${result}['result']    Not the supply address
    #Should Be Equal As Strings    ${jsonRes['result']['info']['contract_invoke']['error_message'] }    Chaincode Error:{\"Error\":\"Not the supply address\"}
