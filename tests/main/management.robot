*** Settings ***
Resource    ../resources/common.robot

Suite Setup    Suite Setup
Test Teardown    Collect Logs

*** Variables ***
${PROJECT_TARBALL}    ${CURDIR}/../testdata/nodered-demo__main@a6293b6.tar.gz

*** Test Cases ***

Install node-red from github url
    ${binary_url}=    Cumulocity.Create Inventory Binary    nodered-demo    nodered-project    file=${CURDIR}/../testdata/nodered-demo.cfg
    ${operation}=    Cumulocity.Install Software
    ...    {"name":"nodered-demo", "version":"latest", "softwareType":"nodered", "url":"${binary_url}"}
    ...    {"name":"active-project", "version":"nodered-demo", "softwareType":"nodered"}
    Operation Should Be SUCCESSFUL    ${operation}
    Cumulocity.Device Should Have Installed Software    {"name": "nodered-demo", "softwareType": "nodered"}
    Cumulocity.Should Have Services    name=nodered-temperature-flow    status=up    service_type=nodered

Install node-red from tarball
    ${binary_url}=    Cumulocity.Create Inventory Binary    nodered-demo    nodered-project    file=${PROJECT_TARBALL}
    ${operation}=    Cumulocity.Install Software
    ...    {"name":"nodered-demo", "version":"latest", "softwareType":"nodered", "url":"${binary_url}"}
    ...    {"name":"active-project", "version":"nodered-demo", "softwareType":"nodered"}
    Operation Should Be SUCCESSFUL    ${operation}
    Cumulocity.Device Should Have Installed Software    {"name": "nodered-demo", "version": "0.0.1", "softwareType": "nodered"}
    Cumulocity.Should Have Services    name=nodered-temperature-flow    status=up    service_type=nodered

Uninstall node-red project via Cumulocity
    ${operation}=    Cumulocity.Uninstall Software    {"name": "nodered-demo", "version": "latest", "softwareType": "nodered"}
    Operation Should Be SUCCESSFUL    ${operation}
    Cumulocity.Device Should Not Have Installed Software    {"name": "nodered-demo", "softwareType": "nodered"}
    Cumulocity.Should Have Services    name=nodered-temperature-flow    status=down    service_type=nodered

Uninstall node-red project via Cumulocity using the active project
    # install first
    ${binary_url}=    Cumulocity.Create Inventory Binary    nodered-demo    nodered-project    file=${PROJECT_TARBALL}
    ${operation}=    Cumulocity.Install Software
    ...    {"name":"nodered-demo", "version":"latest", "softwareType":"nodered", "url":"${binary_url}"}
    ...    {"name":"active-project", "version":"nodered-demo", "softwareType":"nodered"}
    Operation Should Be SUCCESSFUL    ${operation}
    Cumulocity.Device Should Have Installed Software    {"name": "nodered-demo", "version": "0.0.1", "softwareType": "nodered"}
    Cumulocity.Should Have Services    name=nodered-temperature-flow    status=up    service_type=nodered

    # then remove
    ${operation}=    Cumulocity.Uninstall Software    {"name": "active-project", "version": "nodered-demo", "softwareType": "nodered"}
    Operation Should Be SUCCESSFUL    ${operation}
    ${mo}=    Cumulocity.Device Should Have Fragments    c8y_SoftwareList
    Log    ${mo}
    Should Not Contain    ${mo}    nodered-demo
    Should Not Contain    ${mo}    active-project
    Cumulocity.Should Have Services    name=nodered-temperature-flow    status=down    service_type=nodered

Install new project when nodered is not running
    ${operation}=    Cumulocity.Execute Shell Command    text=sudo systemctl stop nodered
    Operation Should Be SUCCESSFUL    ${operation}

    ${binary_url}=    Cumulocity.Create Inventory Binary    nodered-demo    nodered-project    file=${PROJECT_TARBALL}
    ${operation}=    Cumulocity.Install Software
    ...    {"name":"nodered-demo", "version":"latest", "softwareType":"nodered", "url":"${binary_url}"}
    ...    {"name":"active-project", "version":"nodered-demo", "softwareType":"nodered"}
    Operation Should Be SUCCESSFUL    ${operation}
    Cumulocity.Device Should Have Installed Software    {"name": "nodered-demo", "version": "0.0.1", "softwareType": "nodered"}
    Cumulocity.Should Have Services    name=nodered-temperature-flow    status=up    service_type=nodered

Install node-red from github url with space in its name
    [Teardown]    Remove nodered project    nodered demo project
    ${binary_url}=    Cumulocity.Create Inventory Binary    nodered-demo    nodered-project    file=${CURDIR}/../testdata/nodered-demo.cfg
    ${operation}=    Cumulocity.Install Software
    ...    {"name":"nodered demo project", "version":"latest", "softwareType":"nodered", "url":"${binary_url}"}
    ...    {"name":"active-project", "version":"nodered demo project", "softwareType":"nodered"}
    Operation Should Be SUCCESSFUL    ${operation}
    Cumulocity.Device Should Have Installed Software    {"name": "nodered demo project", "softwareType": "nodered"}
    Cumulocity.Should Have Services    name=nodered-temperature-flow    status=up    service_type=nodered

*** Keywords ***

Suite Setup
    ${DEVICE_SN}=    Setup
    Set Suite Variable    $DEVICE_SN
    Cumulocity.External Identity Should Exist    ${DEVICE_SN}
    
Remove nodered project
    [Arguments]    ${name}
    ${operation}=    Cumulocity.Uninstall Software    {"name": "${name}", "version": "latest", "softwareType": "nodered"}
    Operation Should Be SUCCESSFUL    ${operation}
