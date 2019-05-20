*** Settings ***
Library           ../../utilFunc/createToken.py
Resource          ../../utilKwd/utilVariables.txt
Resource          ../../utilKwd/normalKwd.txt
Resource          ../../utilKwd/utilDefined.txt
Resource          ../../utilKwd/behaveKwd.txt

*** Variables ***

*** Test Cases ***
Feature: Create Token
    [Documentation]    Scenario: Verify Reciever's PTN
    ${PTN1}    ${result1}    Given Request getbalance before create token
    ${ret}    When Create token of vote contract
    ${PTNGAIN}    And Calculate gain of recieverAdd    ${PTN1}
    ${PTN2}    ${result2}    And Request getbalance after create token
    Then Assert gain of reciever    ${PTN2}    ${PTNGAIN}

*** Keywords ***
Request getbalance before create token
    personalUnlockAccount    ${geneAdd}
    ${PTN1}    ${result1}    normalGetBalance    ${recieverAdd}
    [Return]    ${PTN1}    ${result1}

Create token of vote contract
    ${geneAdd}    getGeneAdd    ${host}
    ${ccTokenList}    Create List    ${crtTokenMethod}    ${note}    ${tokenDecimal}    ${tokenAmount}    ${voteTime}
    ...    ${commonVoteInfo}
    ${ccList}    Create List    ${geneAdd}    ${recieverAdd}    ${PTNAmount}    ${PTNPoundage}    ${voteContractId}
    ...    ${ccTokenList}    ${pwd}    ${duration}    ${EMPTY}
    ${resp}    setPostRequest    ${host}    ${invokePsMethod}    ${ccList}
    log    ${resp.content}
    Should Contain    ${resp.content}['jsonrpc']    "2.0"    msg="jsonrpc:failed"
    Should Contain    ${resp.content}['id']    1    msg="id:failed"
    ${ret}    Should Match Regexp    ${resp.content}['result']    ${commonResultCode}    msg="result:does't match Result expression"
    [Return]    ${ret}

Calculate gain of recieverAdd
    [Arguments]    ${PTN1}
    ${gain1}    countRecieverPTN    ${PTNAmount}
    ${PTNGAIN}    Evaluate    decimal.Decimal('${PTN1}')+decimal.Decimal('${gain1}')    decimal
    sleep    4
    [Return]    ${PTNGAIN}

Request getbalance after create token
    ${PTN2}    ${result2}    normalGetBalance    ${recieverAdd}
    sleep    3
    [Return]    ${PTN2}    ${result2}

Assert gain of reciever
    [Arguments]    ${PTN2}    ${PTNGAIN}
    Should Be Equal As Numbers    ${PTN2}    ${PTNGAIN}
