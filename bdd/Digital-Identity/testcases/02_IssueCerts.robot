*** Settings ***
Resource          ../../commonlib/pubVariables.robot
Resource          ../../commonlib/pubFuncs.robot
Library           Collections
Library           BuiltIn

*** Test Cases ***
CAIssueIntermedate
    Given CA unlock its account succeed
    ${reqId}=    When CA issues intermediate certificate name cert1 to power succeed
    And Wait for unit about contract to be confirmed by unit height    ${reqId}    ${true}
    Then Power can query his certificate in db

PowerIssueUserCert
    Given Power unlock its account succeed
    ${reqId}=    When Power issues certificate for user succeed
    And Wait for unit about contract to be confirmed by unit height    ${reqId}    ${true}
    Then User can query his certificate in db

*** Keywords ***
Power unlock its account succeed
    ${respJson}=    unlockAccount    ${powerCertHolder}
    Dictionary Should Contain Key    ${respJson}    result
    Should Be Equal    ${respJson["result"]}    ${true}

Power issues certificate for user succeed
    ${params}=    Create List    ${powerCertHolder}    ${userCertHolder}    1    1k    palletone
    ...    user    gptn.mediator1
    ${respJson}=    sendRpcPost    ${host}    wallet_genCert    ${params}    GenCert
    Dictionary Should Contain Key    ${respJson}    result
    ${result}=    Get From Dictionary    ${respJson}    result
    ${reqId}=    Get From Dictionary    ${result}    reqId
    [Return]    ${reqId}

User can query his certificate in db
    ${args}=    Create List    ${getHolderCertMethod}    ${userCertHolder}
    ${params}=    Create List    ${certContractAddr}    ${args}    ${0}
    ${respJson}=    sendRpcPost    ${host}    ${ccqueryMethod}    ${params}    queryCert
    Dictionary Should Contain Key    ${respJson}    result
    ${resultDict}=    Evaluate    ${respJson["result"]}
    Dictionary Should Contain Key    ${resultDict}    MemberCertIDs
    Length Should Be    ${resultDict['MemberCertIDs']}    1
    Dictionary Should Contain Key    ${resultDict['MemberCertIDs'][0]}    CertID
    ${CertID}=    Evaluate    ${resultDict}['MemberCertIDs'][0]['CertID']
    Set Global Variable    ${userCertID}    ${CertID}

CA unlock its account succeed
    ${respJson}=    unlockAccount    ${caCertHolder}
    Dictionary Should Contain Key    ${respJson}    result
    Should Be Equal    ${respJson["result"]}    ${true}

CA issues intermediate certificate name cert1 to power succeed
    ${args}=    Create List    addServerCert    ${powerCertBytes}
    ${params}=    genInvoketxParams    ${caCertHolder}    ${caCertHolder}    100    100    ${certContractAddr}
    ...    ${args}    ${null}
    ${respJson}=    sendRpcPost    ${host}    ${ccinvokeMethod}    ${params}    addServerCert
    Dictionary Should Contain Key    ${respJson}    result
    ${result}=    Get From Dictionary    ${respJson}    result
    ${reqId}=    Get From Dictionary    ${result}    reqId
    [Return]    ${reqId}

Power can query his certificate in db
    ${args}=    Create List    ${getHolderCertMethod}    ${powerCertHolder}
    ${params}=    Create List    ${certContractAddr}    ${args}    ${0}
    ${respJson}=    sendRpcPost    ${host}    ${ccqueryMethod}    ${params}    queryCert
    Dictionary Should Contain Key    ${respJson}    result
    ${resultDict}=    Evaluate    ${respJson["result"]}
    Dictionary Should Contain Key    ${resultDict}    IntermediateCertIDs
    Length Should Be    ${resultDict['IntermediateCertIDs']}    1
    Dictionary Should Contain Key    ${resultDict['IntermediateCertIDs'][0]}    CertID
    ${CertID}=    Evaluate    ${resultDict}['IntermediateCertIDs'][0]['CertID']
    Set Global Variable    ${powerCertID}    ${CertID}
