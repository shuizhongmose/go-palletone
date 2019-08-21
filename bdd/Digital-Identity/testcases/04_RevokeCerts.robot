*** Settings ***
Library           Collections
Library           BuiltIn
Library           DateTime
Library           String
Resource          ../../commonlib/pubVariables.robot
Resource          ../../commonlib/pubFuncs.robot

*** Test Cases ***
PowerRevokeUserCert
    Given Power unlock his account succeed
    ${reqId}=    When Power revoke user certificate succeed
    And Wait for unit about contract to be confirmed by unit height    ${reqId}    ${true}
    And Sleep    10
    Then Power can query his issued CRL file
    And User certificate revocation time is before now

CARevokePowerCert
    Given CA unlock his account succeed
    ${reqId}=    When CA revoke power certificate succeed
    And Wait for unit about contract to be confirmed by unit height    ${reqId}    ${true}
    And Sleep    10
    Then CA can query his issued CRL file
    And Power certificate revocation time is before now

*** Keywords ***
Power unlock his account succeed
    ${respJson}=    unlockAccount    ${powerCertHolder}
    Dictionary Should Contain Key    ${respJson}    result
    Should Be Equal    ${respJson["result"]}    ${true}

Power revoke user certificate succeed
    ${params}=    Create List    ${powerCertHolder}    1    ${userCertHolder}
    ${respJson}=    sendRpcPost    ${host}    wallet_revokeCert    ${params}    RevokeCert
    Dictionary Should Contain Key    ${respJson}    result
    ${result}=    Get From Dictionary    ${respJson}    result
    ${reqId}=    Get From Dictionary    ${result}    reqId
    [Return]    ${reqId}

Power can query his issued CRL file
    ${args}=    Create List    ${queryCRLMethod}    ${powerCertHolder}
    ${params}=    Create List    ${certContractAddr}    ${args}    ${0}
    ${respJson}=    sendRpcPost    ${host}    ${ccqueryMethod}    ${params}    queryCRL
    Dictionary Should Contain Key    ${respJson}    result
    ${bytes}=    Evaluate    ${respJson['result']}
    Length Should Be    ${bytes}    1

User certificate revocation time is before now
    ${args}=    Create List    ${getHolderCertMethod}    ${userCertHolder}
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
    ${params}=    Create List    ${tokenHolder}    1    ${powerCertHolder}
    ${respJson}=    sendRpcPost    ${host}    wallet_revokeCert    ${params}    RevokeCert
    Dictionary Should Contain Key    ${respJson}    result
    ${result}=    Get From Dictionary    ${respJson}    result
    ${reqId}=    Get From Dictionary    ${result}    reqId
    [Return]    ${reqId}

CA can query his issued CRL file
    ${args}=    Create List    ${queryCRLMethod}    ${tokenHolder}
    ${params}=    Create List    ${certContractAddr}    ${args}    ${0}
    ${respJson}=    sendRpcPost    ${host}    ${ccqueryMethod}    ${params}    queryCRL
    Dictionary Should Contain Key    ${respJson}    result
    ${bytes}=    Evaluate    ${respJson['result']}
    Length Should Be    ${bytes}    1

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
