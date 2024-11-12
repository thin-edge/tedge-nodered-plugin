*** Settings ***
Resource    ../resources/common.robot

Suite Setup    Custom Setup
Test Teardown    Collect Logs
Suite Teardown    Collect Logs

*** Test Cases ***

node-red project should be running and processing thin-edge telemetry
    ${test_start}=    DeviceLibrary.Get Test Start Time

    # Send an initial measurement to reset the flow's state
    ${operation}=    Cumulocity.Execute Shell Command    text=tedge mqtt pub 'te/device/main///m/env' '{"temperature":10}'
    Cumulocity.Operation Should Be SUCCESSFUL    ${operation}

    # Send a big enough temperature spike to invoke the flow's logic
    ${operation}=    Cumulocity.Execute Shell Command    text=tedge mqtt pub 'te/device/main///m/env' '{"temperature":50}'
    Cumulocity.Operation Should Be SUCCESSFUL    ${operation}
    Cumulocity.Device Should Have Event/s     type=temperatureChange   expected_text=Temperature changed by ≥10°C. new_value=50°C    after=${test_start}

node-red status should publish to health endpoint
    Cumulocity.Should Have Services    name=nodered-temperature-flow    service_type=nodered    status=up

*** Keywords ***

Custom Setup
    ${DEVICE_SN}=    Setup
    Set Suite Variable    $DEVICE_SN
    Cumulocity.External Identity Should Exist    ${DEVICE_SN}

    # Install
    # TODO: Create conditional installation
    ${binary_url}=    Cumulocity.Create Inventory Binary    nodered    application/x-yaml    file=${CURDIR}/../testdata/docker-compose.nodered-project.yaml
    ${operation}=    Cumulocity.Install Software
    ...    {"name":"nodered", "version":"1.0.0", "softwareType":"container-group", "url":"${binary_url}"}
    Operation Should Be SUCCESSFUL    ${operation}

    ${binary_url}=    Cumulocity.Create Inventory Binary    nodered-demo    nodered-project    file=${CURDIR}/../testdata/nodered-demo.json
    ${operation}=    Cumulocity.Install Software
    ...    {"name":"nodered-demo", "version":"latest", "softwareType":"nodered-project", "url":"${binary_url}"}
    Operation Should Be SUCCESSFUL    ${operation}    timeout=90
