*** Settings ***
Default Tags      normal
Library           ../../utilFunc/createToken.py
Resource          ../../utilKwd/utilVariables.txt
Resource          ../../utilKwd/normalKwd.txt
Resource          ../../utilKwd/utilDefined.txt
Resource          ../../utilKwd/behaveKwd.txt

*** Variables ***
${preTokenId}     QA054

*** Test Cases ***
Scenario: 20Contract - Transfer Token
    [Documentation]    Verify Reciever's Token
    Given Request getbalance before create token
    ${ret}    When Request normal CcinvokePass
    ${key}    ${item}    And Request getbalance after create token
    Then Assert gain    ${key}    ${item}

*** Keywords ***
Request getbalance before create token
    ${geneAdd}    getGeneAdd    ${host}
    Set Suite Variable    ${geneAdd}    ${geneAdd}
<<<<<<< HEAD
    sleep    4
=======
>>>>>>> master

Request normal CcinvokePass
    ${ccList}    Create List    ${crtTokenMethod}    ${evidence}    ${preTokenId}    ${tokenDecimal}    ${tokenAmount}
    ...    ${geneAdd}
    ${ret}    normalCcinvokePass    ${commonResultCode}    ${geneAdd}    ${recieverAdd}    ${PTNAmount}    ${PTNPoundage}
    ...    ${20ContractId}    ${ccList}
    [Return]    ${ret}

Request getbalance after create token
    sleep    5
    ${PTN2}    ${result2}    normalGetBalance    ${geneAdd}
<<<<<<< HEAD
    sleep    2
=======
>>>>>>> master
    ${key}    getTokenId    ${preTokenId}    ${result2['result']}
    ${item}    Set Variable    0
    ${tokenResult}    transferToken    ${key}    ${geneAdd}    ${recieverAdd}    ${gain}    ${PTNPoundage}
    ...    ${evidence}    ${duration}
    sleep    4
    [Return]    ${key}    ${item}

Assert gain
    [Arguments]    ${key}    ${item}
    ${item1}    Evaluate    ${item}+${gain}
    sleep    4
    ${RecPTN2}    ${RecResult2}    normalGetBalance    ${recieverAdd}
<<<<<<< HEAD
    sleep    2
=======
>>>>>>> master
    ${item2}    Get From Dictionary    ${RecResult2['result']}    ${key}
    Should Be Equal As Numbers    ${item2}    ${item1}
