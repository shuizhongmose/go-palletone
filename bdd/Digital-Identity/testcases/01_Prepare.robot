*** Settings ***
Resource          ../../commonlib/pubVariables.robot
Resource          ../../commonlib/pubFuncs.robot
Resource          ../../commonlib/setups.robot
Library           Collections
Library           BuiltIn
Library           OperatingSystem
Library           String

*** Variables ***
#${CertFilePath}    C:/Users/Administrator/Desktop/tmp/digital-identity
#${CrlFilePath}    C:/Users/Administrator/Desktop/tmp/digital-identity
${CertFilePath}    ~/cawork/immediateca/
${CrlFilePath}    ~/cawork/root/msp/crls/

*** Test Cases ***
testprepare
    queryTokenHolder    ${false}
    queryCACertID
    # new account
    ${user}=    newAccount
    Set Global Variable    ${powerCertHolder}    ${user}
    ${user}=    newAccount
    Set Global Variable    ${userCertHolder}    ${user}
    ${user}=    newAccount
    Set Global Variable    ${userCertHolder2}    ${user}
    # transfer ptn to power and user
    transferPtnTo    ${powerCertHolder}    10000
    transferPtnTo    ${userCertHolder}    10000
    transferPtnTo    ${userCertHolder2}    10000
    # query power cert bytes from ~/cawork/immediateca/
    ${cert}=    Get File    ${CertFilePath}/ca-cert.pem
    Set Global Variable    ${powerCertBytes}    ${cert}
    # query immediate crl bytes
    ${crl}=    Get File    ${CrlFilePath}/crl.pem
    Set Global Variable    ${immediateCrlBytes}    ${crl}
