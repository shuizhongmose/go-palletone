*** Settings ***
Library           Collections
Library           BuiltIn
Library           DateTime
Library           String
Resource          ../../commonlib/pubVariables.robot
Resource          ../../commonlib/pubFuncs.robot

*** Test Cases ***
PowerRevokeUser1Cert1
    Given Power unlock his account succeed
    ${reqId}=    When Power revoke user certificate succeed    ${userCertHolder}
    And Wait for unit about contract to be confirmed by unit height    ${reqId}    ${true}
    ${issuer}=    And Get invoke payload info    ${reqId}
    Then Power can query his issued CRL file    ${issuer}
    And User certificate revocation time is before now    ${userCertHolder}

CARevokePowerCert
    Given CA unlock his account succeed
    ${reqId}=    When CA revoke power certificate succeed
    And Wait for unit about contract to be confirmed by unit height    ${reqId}    ${true}
    ${issuer}=    And Get invoke payload info    ${reqId}
    Then CA can query his issued CRL file    ${issuer}
    And User certificate revocation time is before now    ${userCertHolder2}

*** Keywords ***
Power unlock his account succeed
    ${respJson}=    unlockAccount    ${powerCertHolder}
    Dictionary Should Contain Key    ${respJson}    result
    Should Be Equal    ${respJson["result"]}    ${true}

Power revoke user certificate succeed
    [Arguments]    ${userAddr}
    ${params}=    Create List    ${powerCertHolder}    1    ${userAddr}
    ${respJson}=    sendRpcPost    ${host}    wallet_revokeCert    ${params}    RevokeCert
    Dictionary Should Contain Key    ${respJson}    result
    ${result}=    Get From Dictionary    ${respJson}    result
    ${reqId}=    Get From Dictionary    ${result}    reqId
    [Return]    ${reqId}

Power can query his issued CRL file
    [Arguments]    ${issuer}
    ${args}=    Create List    ${queryCRLMethod}    ${issuer}
    ${params}=    Create List    ${certContractAddr}    ${args}    ${0}
    ${respJson}=    sendRpcPost    ${host}    ${ccqueryMethod}    ${params}    queryCRL
    Dictionary Should Contain Key    ${respJson}    result
    ${bytes}=    Evaluate    ${respJson['result']}
    Length Should Be    ${bytes}    1

User certificate revocation time is before now
    [Arguments]    ${userAddr}
    ${args}=    Create List    ${getHolderCertMethod}    ${userAddr}
    ${params}=    Create List    ${certContractAddr}    ${args}    ${0}
    ${respJson}=    sendRpcPost    ${host}    ${ccqueryMethod}    ${params}    queryCert
    Dictionary Should Contain Key    ${respJson}    result
    ${resultDict}=    Evaluate    ${respJson["result"]}
    Dictionary Should Contain Key    ${resultDict}    MemberCertIDs
    Length Should Be    ${resultDict['MemberCertIDs']}    1
    Dictionary Should Contain Key    ${resultDict['MemberCertIDs'][0]}    CertID
    ${now}=    Get Current Date    UTC
    ${words}=    Split String    ${resultDict['MemberCertIDs'][0]['RecovationTime']}    ${SPACE}
    Length Should Be    ${words}    4
    ${sRevocationTime}=    catenate    ${words[0]}    ${words[1]}
    ${sRevocationTime}=    catenate    SEPARATOR=    ${sRevocationTime}    .000
    ${revocationTime}=    Convert Date    ${sRevocationTime}
    Run Keyword If    '${now}'>'${revocationTime}'    log    1
    ...    ELSE    Fail    Power invoke user certificate failed

CA unlock his account succeed
    ${respJson}=    unlockAccount    ${tokenHolder}
    Dictionary Should Contain Key    ${respJson}    result
    Should Be Equal    ${respJson["result"]}    ${true}

CA revoke power certificate succeed
    ${args}=    Create List    ${addCRLMethod}    ${immediateCrlBytes}
    ${respJson}=    invokeContract    ${tokenHolder}    ${tokenHolder}    100    100    ${certContractAddr}
    ...    ${args}
    Dictionary Should Contain Key    ${respJson}    result
    ${result}=    Get From Dictionary    ${respJson}    result
    ${reqId}=    Get From Dictionary    ${result}    reqId
    [Return]    ${reqId}

CA can query his issued CRL file
    [Arguments]    ${issuer}
    ${args}=    Create List    ${queryCRLMethod}    ${issuer}
    ${params}=    Create List    ${certContractAddr}    ${args}    ${0}
    ${respJson}=    sendRpcPost    ${host}    ${ccqueryMethod}    ${params}    queryCRL
    Dictionary Should Contain Key    ${respJson}    result
    ${bytes}=    Evaluate    ${respJson['result']}
    Should Be Equal    ${bytes}    ${immediateCrlBytes}

Power certificate revocation time is before now
    ${args}=    Create List    ${getHolderCertMethod}    ${powerCertHolder}
    ${params}=    Create List    ${certContractAddr}    ${args}    ${0}
    ${respJson}=    sendRpcPost    ${host}    ${ccqueryMethod}    ${params}    queryCert
    Dictionary Should Contain Key    ${respJson}    result
    ${resultDict}=    Evaluate    ${respJson["result"]}
    Dictionary Should Contain Key    ${resultDict}    IntermediateCertIDs
    Length Should Be    ${resultDict['IntermediateCertIDs']}    1
    Dictionary Should Contain Key    ${resultDict['IntermediateCertIDs'][0]}    CertID
    ${now}=    Get Current Date    UTC
    ${words}=    Split String    ${resultDict['IntermediateCertIDs'][0]['RecovationTime']}    ${SPACE}
    Length Should Be    ${words}    4
    ${sRevocationTime}=    catenate    ${words[0]}    ${words[1]}
    ${sRevocationTime}=    catenate    SEPARATOR=    ${sRevocationTime}    .000
    ${revocationTime}=    Convert Date    ${sRevocationTime}
    Run Keyword If    '${now}'>'${revocationTime}'    log    1
    ...    ELSE    Fail    ca invoke power certificate failed
