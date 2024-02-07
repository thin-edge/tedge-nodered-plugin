*** Settings ***
Resource    ../resources/common.robot
Library    Cumulocity
Library    DeviceLibrary

Suite Setup    Set Main Device

*** Variables ***
${PROJECT_TARBALL}    ${CURDIR}/../testdata/nodered-demo__main@a6293b6.tar.gz

*** Test Cases ***

Install node-red from github url
    ${binary_url}=    Cumulocity.Create Inventory Binary    nodered-demo    nodered-project    file=${CURDIR}/../testdata/nodered-demo.cfg
    ${operation}=    Cumulocity.Install Software    nodered-demo,latest::nodered,${binary_url}    active-project,nodered-demo::nodered
    Operation Should Be SUCCESSFUL    ${operation}
    Cumulocity.Device Should Have Installed Software    nodered-demo
    Cumulocity.Should Have Services    name=nodered-temperature-flow    status=up    service_type=nodered

Install node-red from tarball
    ${binary_url}=    Cumulocity.Create Inventory Binary    nodered-demo    nodered-project    file=${PROJECT_TARBALL}
    ${operation}=    Cumulocity.Install Software    nodered-demo,latest::nodered,${binary_url}    active-project,nodered-demo::nodered
    Operation Should Be SUCCESSFUL    ${operation}
    Cumulocity.Device Should Have Installed Software    nodered-demo,0.0.1
    Cumulocity.Should Have Services    name=nodered-temperature-flow    status=up    service_type=nodered

Uninstall node-red project via Cumulocity
    ${operation}=    Cumulocity.Create Operation    fragments={"c8y_SoftwareUpdate":[{"name":"nodered-demo","version":"latest::nodered","url":"","action":"delete"}]}    description=Remove nodered-demo package
    Operation Should Be SUCCESSFUL    ${operation}
    ${mo}=    Cumulocity.Device Should Have Fragments    c8y_SoftwareList
    Log    ${mo}
    Should Not Contain    ${mo}    nodered-demo
    Cumulocity.Should Have Services    name=nodered-temperature-flow    status=down    service_type=nodered

Uninstall node-red project via Cumulocity using the active project
    # install first
    ${binary_url}=    Cumulocity.Create Inventory Binary    nodered-demo    nodered-project    file=${PROJECT_TARBALL}
    ${operation}=    Cumulocity.Install Software    nodered-demo,latest::nodered,${binary_url}    active-project,nodered-demo::nodered
    Operation Should Be SUCCESSFUL    ${operation}
    Cumulocity.Device Should Have Installed Software    nodered-demo,0.0.1
    Cumulocity.Should Have Services    name=nodered-temperature-flow    status=up    service_type=nodered

    # then remove
    ${operation}=    Cumulocity.Create Operation    fragments={"c8y_SoftwareUpdate":[{"name":"active-project","version":"nodered-demo::nodered","url":"","action":"delete"}]}    description=Remove nodered-demo package
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
    ${operation}=    Cumulocity.Install Software    nodered-demo,latest::nodered,${binary_url}    active-project,nodered-demo::nodered
    Operation Should Be SUCCESSFUL    ${operation}
    Cumulocity.Device Should Have Installed Software    nodered-demo,0.0.1
    Cumulocity.Should Have Services    name=nodered-temperature-flow    status=up    service_type=nodered

Install node-red from github url with space in its name
    [Teardown]    Remove nodered project    nodered demo project
    ${binary_url}=    Cumulocity.Create Inventory Binary    nodered-demo    nodered-project    file=${CURDIR}/../testdata/nodered-demo.cfg
    ${operation}=    Cumulocity.Install Software    nodered demo project,latest::nodered,${binary_url}    active-project,nodered demo project::nodered
    Operation Should Be SUCCESSFUL    ${operation}
    Cumulocity.Device Should Have Installed Software    nodered demo project
    Cumulocity.Should Have Services    name=nodered-temperature-flow    status=up    service_type=nodered

*** Keywords ***

Remove nodered project
    [Arguments]    ${name}
    ${operation}=    Cumulocity.Create Operation    fragments={"c8y_SoftwareUpdate":[{"name":"${name}","version":"latest::nodered","url":"","action":"delete"}]}    description=Remove 'nodered demo project' package
    Operation Should Be SUCCESSFUL    ${operation}
