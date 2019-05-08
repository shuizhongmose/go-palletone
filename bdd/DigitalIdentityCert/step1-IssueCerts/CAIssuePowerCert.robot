*** Settings ***
Test Setup        beforeIssueCerts
Library           Collections
Resource          ../pubVariables.robot
Resource          ../pubFuncs.robot
Resource          ../setups.robot

*** Test Cases ***
CAIssueIntermedate
    Given ca certificate exists
    When ca unlock its account succeed
    and ca issues intermediate certificate name cert1 to power succeed
    and wait for transaction being packaged
    Then power can query his certificate in db

*** Keywords ***
ca certificate exists
    Log    "ca certificate exists"
    ${args}=    Create List    getRootCAHoler
    ${params}=    Create List    ${certContractAddr}    ${args}
    # send post
    ${respJson}=    sendRpcPost    ${queryMethod}    ${params}    queryCAHolder
    # check result
    Dictionary Should Contain Key    ${respJson}    result
    Should Be Equal    ${respJson["result"]}    ${caCertHolder}

ca unlock its account succeed
    Log    "ca unlock its account succeed"
    ${respJson}=    unlockAccount    ${caCertHolder}
    Dictionary Should Contain Key    ${respJson}    result
    Should Be Equal    ${respJson["result"]}    ${true}

ca issues intermediate certificate name cert1 to power succeed
    Log    "ca issues intermediate certificate name cert1 to user1 succeed"
    ${args}=    Create List    addServerCert    ${powerCertHolder}    ${powerCertBytes}
    ${params}=    Create List    ${caCertHolder}    ${caCertHolder}    1    1    ${certContractAddr}
    ...    ${args}    ${null}
    ${respJson}=    sendRpcPost    ${invokeMethod}    ${params}    addServerCert
    Dictionary Should Contain Key    ${respJson}    result

power can query his certificate in db
    ${args}=    Create List    ${getHolderCertMethod}    ${powerCertHolder}
    ${params}=    Create List    ${certContractAddr}    ${args}
    ${respJson}=    sendRpcPost    ${queryMethod}    ${params}    queryCert
    Dictionary Should Contain Key    ${respJson}    result
    ${resultDict}=    Evaluate    ${respJson["result"]}
    Dictionary Should Contain Key    ${resultDict}    IntermediateCertIDs
    Length Should Be    ${resultDict['IntermediateCertIDs']}    1
    Dictionary Should Contain Key    ${resultDict['IntermediateCertIDs'][0]}    CertID
    Should Be Equal    ${resultDict['IntermediateCertIDs'][0]['CertID']}    ${powerCertID}
