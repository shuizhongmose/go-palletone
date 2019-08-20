*** Settings ***
Library           Collections
Library           BuiltIn
Resource          ../../commonlib/pubVariables.robot
Resource          ../../commonlib/pubFuncs.robot

*** Test Cases ***
#CAUseCert
#    Given CA unlock account succed
#    ${reqId}=    When CA uses debug contract to test getRequesterCert without error
#    And Wait for unit about contract to be confirmed by unit height    ${reqId}     ${true}
#    ${reqId}=    Then CA uses debug contract to test checkRequesterCert without error
#    And Wait for unit about contract to be confirmed by unit height    ${reqId}     ${true}
#
#PowerUseCert
#    Given Power unlock account succed
#    ${reqId}=    When Power uses debug contract to test getRequesterCert without error
#    And Wait for unit about contract to be confirmed by unit height    ${reqId}     ${true}
#    ${reqId}=    Then Power uses debug contract to test checkRequesterCert without error
#    And Wait for unit about contract to be confirmed by unit height    ${reqId}     ${true}

UserUseCert
    Given User unlock account succed
    ${reqId}=    When User uses debug contract to test getRequesterCert without error
    And Wait for unit about contract to be confirmed by unit height    ${reqId}     ${true}
    ${reqId}=    Then User uses debug contract to test checkRequesterCert without error
    And Wait for unit about contract to be confirmed by unit height    ${reqId}     ${true}

*** Keywords ***
CA unlock account succed
    Log    "CA unlock account succed"
    ${respJson}=    unlockAccount    ${caCertHolder}
    Dictionary Should Contain Key    ${respJson}    result
    Should Be Equal    ${respJson["result"]}    ${true}

CA uses debug contract to test getRequesterCert without error
    ${args}=    Create List    ${getRequesterCertMethod}
    ${params}=    genInvoketxParams    ${caCertHolder}    ${caCertHolder}    1    1    ${debugContractAddr}
    ...    ${args}    ${caCertID}
    ${respJson}=    sendRpcPost     ${host}       ${ccinvokeMethod}    ${params}    getRequesterCert
    Dictionary Should Contain Key    ${respJson}    result
    ${result}=    Get From Dictionary    ${respJson}    result
    ${reqId}=    Get From Dictionary    ${result}    reqId
    [Return]    ${reqId}

CA uses debug contract to test checkRequesterCert without error
    ${args}=    Create List    ${checkRequesterCertMethod}
    ${params}=    genInvoketxParams    ${caCertHolder}    ${caCertHolder}    1    1    ${debugContractAddr}
    ...    ${args}    ${caCertID}
    ${respJson}=    sendRpcPost     ${host}       ${ccinvokeMethod}    ${params}    checkRequesterCert
    Dictionary Should Contain Key    ${respJson}    result
    ${result}=    Get From Dictionary    ${respJson}    result
    ${reqId}=    Get From Dictionary    ${result}    reqId
    [Return]    ${reqId}

Power unlock account succed
    ${respJson}=    unlockAccount    ${powerCertHolder}
    Dictionary Should Contain Key    ${respJson}    result
    Should Be Equal    ${respJson["result"]}    ${true}

Power uses debug contract to test getRequesterCert without error
    ${args}=    Create List    ${getRequesterCertMethod}
    ${params}=    genInvoketxParams    ${powerCertHolder}    ${powerCertHolder}    1    1    ${debugContractAddr}
    ...    ${args}    ${powerCertID}
    ${respJson}=    sendRpcPost      ${host}      ${ccinvokeMethod}    ${params}    getRequesterCert
    Dictionary Should Contain Key    ${respJson}    result
    ${result}=    Get From Dictionary    ${respJson}    result
    ${reqId}=    Get From Dictionary    ${result}    reqId
    [Return]    ${reqId}

Power uses debug contract to test checkRequesterCert without error
    ${args}=    Create List    ${checkRequesterCertMethod}
    ${params}=    genInvoketxParams    ${powerCertHolder}    ${powerCertHolder}    1    1    ${debugContractAddr}
    ...    ${args}    ${powerCertID}
    ${respJson}=    sendRpcPost     ${host}       ${ccinvokeMethod}    ${params}    checkRequesterCert
    Dictionary Should Contain Key    ${respJson}    result
    ${result}=    Get From Dictionary    ${respJson}    result
    ${reqId}=    Get From Dictionary    ${result}    reqId
    [Return]    ${reqId}

User unlock account succed
    ${respJson}=    unlockAccount    ${userCertHolder}
    Dictionary Should Contain Key    ${respJson}    result
    Should Be Equal    ${respJson["result"]}    ${true}

User uses debug contract to test getRequesterCert without error
    ${args}=    Create List    ${getRequesterCertMethod}
    ${params}=    genInvoketxParams    ${userCertHolder}    ${userCertHolder}    1    1    ${debugContractAddr}
    ...    ${args}    ${userCertID}
    ${respJson}=    sendRpcPost    ${host}    ${ccinvokeMethod}    ${params}    getRequesterCert
    Dictionary Should Contain Key    ${respJson}    result
    ${result}=    Get From Dictionary    ${respJson}    result
    ${reqId}=    Get From Dictionary    ${result}    reqId
    [Return]    ${reqId}

User uses debug contract to test checkRequesterCert without error
    ${args}=    Create List    ${checkRequesterCertMethod}
    ${params}=    genInvoketxParams    ${userCertHolder}    ${userCertHolder}    1    1    ${debugContractAddr}
    ...    ${args}    ${userCertID}
    ${respJson}=    sendRpcPost     ${host}       ${ccinvokeMethod}    ${params}    checkRequesterCert
    Dictionary Should Contain Key    ${respJson}    result
    ${result}=    Get From Dictionary    ${respJson}    result
    ${reqId}=    Get From Dictionary    ${result}    reqId
    [Return]    ${reqId}

