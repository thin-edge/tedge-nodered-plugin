*** Settings ***
Resource    ../resources/common.robot

Suite Setup    Suite Setup
Test Teardown    Collect Logs

*** Variables ***
${PROJECT_TARBALL}    ${CURDIR}/../testdata/nodered-demo__next@cc38b0a.tar.gz

*** Test Cases ***

Install Flows from Cumulocity
    ${binary_url}=    Cumulocity.Create Inventory Binary    flow1    application/json    file=${CURDIR}/../testdata/flow1.json
    ${operation}=    Cumulocity.Install Software
    ...    {"name":"flow1", "version":"1.0.0", "softwareType":"nodered-flows", "url":"${binary_url}"}
    Operation Should Be SUCCESSFUL    ${operation}
    Cumulocity.Device Should Have Installed Software    {"name": "flow1", "version":"1.0.0", "softwareType":"nodered-flows"}
    Cumulocity.Should Have Services    name=nodered-temperature-flow    status=up    service_type=nodered

Replace existing Flow
    ${binary_url}=    Cumulocity.Create Inventory Binary    nodered-demo    nodered-project    file=${CURDIR}/../testdata/flow2.json
    ${operation}=    Cumulocity.Install Software
    ...    {"name":"flow2", "version":"1.2.3", "softwareType":"nodered-flows", "url":"${binary_url}"}
    Operation Should Be SUCCESSFUL    ${operation}
    Cumulocity.Device Should Have Installed Software    {"name": "flow2", "version":"1.2.3", "softwareType":"nodered-flows"}
    Cumulocity.Device Should Not Have Installed Software    {"name": "flow1","softwareType": "nodered-flows"}
    Cumulocity.Should Have Services    name=nodered-temperature-flow2    status=up    service_type=nodered
    Cumulocity.Should Have Services    name=nodered-temperature-flow    status=down    service_type=nodered

Uninstall Flow
    ${operation}=    Cumulocity.Uninstall Software    {"name": "flow2", "version": "1.2.3", "softwareType": "nodered-flows"}
    Operation Should Be SUCCESSFUL    ${operation}
    Cumulocity.Device Should Not Have Installed Software    {"name": "flow2", "softwareType": "nodered-flows"}
    Cumulocity.Should Have Services    name=nodered-temperature-flow2    status=down    service_type=nodered


*** Keywords ***

Suite Setup
    ${DEVICE_SN}=    Setup
    Set Suite Variable    $DEVICE_SN
    Cumulocity.External Identity Should Exist    ${DEVICE_SN}

    # Install
    ${binary_url}=    Cumulocity.Create Inventory Binary    nodered    application/x-yaml    file=${CURDIR}/../testdata/docker-compose.nodered-flows.yaml
    ${operation}=    Cumulocity.Install Software
    ...    {"name":"nodered", "version":"1.0.0", "softwareType":"container-group", "url":"${binary_url}"}
    Operation Should Be SUCCESSFUL    ${operation}    timeout=90
