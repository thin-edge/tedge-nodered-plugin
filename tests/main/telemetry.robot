*** Settings ***
Resource    ../resources/common.robot
Library    Cumulocity
Library    DeviceLibrary

Suite Setup    Custom Setup

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
    Set Main Device
    ${binary_url}=    Cumulocity.Create Inventory Binary    nodered-demo    nodered-project    file=${CURDIR}/../testdata/nodered-demo.cfg
    ${operation}=    Cumulocity.Install Software    nodered-demo,latest::nodered,${binary_url}    active-project,nodered-demo::nodered
    Operation Should Be SUCCESSFUL    ${operation}
