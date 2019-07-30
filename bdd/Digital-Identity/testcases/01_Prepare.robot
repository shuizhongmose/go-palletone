*** Settings ***
Resource          ../../commonlib/pubVariables.robot
Resource          ../../commonlib/pubFuncs.robot
Resource          ../../commonlib/setups.robot
Library           Collections
Library           BuiltIn

*** Test Cases ***
testprepare
    queryTokenHolder    ${false}
    queryCAHolder
    # new account
    ${user}=    newAccount
    Set Global Variable    ${powerCertHolder}    ${user}
    ${user}=    newAccount
    Set Global Variable    ${userCertHolder}    ${user}
    # transfer ptn to power and user
    transferPtnTo    ${powerCertHolder}    1000
    transferPtnTo    ${userCertHolder}    1000
    # query power cert bytes from ~/cawork/immediateca/
