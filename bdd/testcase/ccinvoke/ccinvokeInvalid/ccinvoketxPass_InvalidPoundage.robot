*** Settings ***
Suite Setup       getlistAccounts
Force Tags        invalidAdd
Default Tags      invalidAdd
Library           RequestsLibrary
Library           Collections
Library           ../../utilFunc/createToken.py
Resource          ../../utilKwd/invalidKwd.txt
Resource          ../../utilKwd/utilDefined.txt
Resource          ../../utilKwd/behaveKwd.txt
Resource          ../../utilKwd/normalKwd.txt

*** Variables ***
${host}           http://localhost:8545/
${method}         contract_ccinvoketsPass

*** Test Cases ***
Scenario: invalidPoundage
    [Template]    InvalidCcinvoke
    200    ${EMPTY}    PCGTta3M4t3yXu8uRgkKvaWd2d8DREThG43    createToken    QA666    evidence    2
    ...    1000    1    ${6000}    ${Empty}    -32000    transaction fee cannot be 0    ${listAccounts[0]}
    ...    ${listAccounts[1]}    ${listAccounts[1]}
    200    0    PCGTta3M4t3yXu8uRgkKvaWd2d8DREThG43    createToken    QA666    evidence    2
    ...    1000    1    ${6000}    ${Empty}    -32000    transaction fee cannot be 0    ${listAccounts[0]}
    ...    ${listAccounts[1]}    ${listAccounts[1]}
    200    -3    PCGTta3M4t3yXu8uRgkKvaWd2d8DREThG43    createToken    QA666    evidence    2
    ...    1000    1    ${6000}    ${Empty}    -32000    transaction fee cannot be 0    ${listAccounts[0]}
    ...    ${listAccounts[1]}    ${listAccounts[1]}
    200    -0.3    PCGTta3M4t3yXu8uRgkKvaWd2d8DREThG43    createToken    QA666    evidence    2
    ...    1000    1    ${6000}    ${Empty}    -32000    transaction fee cannot be 0    ${listAccounts[0]}
    ...    ${listAccounts[1]}    ${listAccounts[1]}
    200    0.5    PCGTta3M4t3yXu8uRgkKvaWd2d8DREThG43    createToken    QA666    evidence    2
    ...    1000    1    ${6000}    ${Empty}    -32000    transaction fee cannot be 0    ${listAccounts[0]}
    ...    ${listAccounts[1]}    ${listAccounts[1]}
    200    a    PCGTta3M4t3yXu8uRgkKvaWd2d8DREThG43    createToken    QA666    evidence    2
    ...    1000    1    ${6000}    ${Empty}    -32000    transaction fee cannot be 0    ${listAccounts[0]}
    ...    ${listAccounts[1]}    ${listAccounts[1]}
    200    ${SPACE}    PCGTta3M4t3yXu8uRgkKvaWd2d8DREThG43    createToken    QA666    evidence    2
    ...    1000    1    ${6000}    ${Empty}    -32000    transaction fee cannot be 0    ${listAccounts[0]}
    ...    ${listAccounts[1]}    ${listAccounts[1]}
    200    $    PCGTta3M4t3yXu8uRgkKvaWd2d8DREThG43    createToken    QA666    evidence    2
    ...    1000    1    ${6000}    ${Empty}    -32000    transaction fee cannot be 0    ${listAccounts[0]}
    ...    ${listAccounts[1]}    ${listAccounts[1]}

*** Keywords ***
