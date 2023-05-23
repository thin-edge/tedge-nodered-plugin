*** Settings ***
Resource    ../resources/common.robot
Library    Cumulocity
Library    DeviceLibrary

Suite Setup    Set Main Device

*** Test Cases ***

Install node-red from github url
    ${binary_url}=    Cumulocity.Create Inventory Binary    nodered-demo    nodered-project    file=${CURDIR}/../testdata/nodered-demo.cfg
    ${operation}=    Cumulocity.Install Software    nodered-demo,latest::nodered,${binary_url}    active-project,nodered-demo::nodered
    Operation Should Be SUCCESSFUL    ${operation}
    Cumulocity.Device Should Have Installed Software    nodered-demo

Install node-red from tarball
    ${binary_url}=    Cumulocity.Create Inventory Binary    nodered-demo    nodered-project    file=${CURDIR}/../testdata/nodered-demo__main@c7c6b5d.tar.gz
    ${operation}=    Cumulocity.Install Software    nodered-demo,latest::nodered,${binary_url}    active-project,nodered-demo::nodered
    Operation Should Be SUCCESSFUL    ${operation}
    Cumulocity.Device Should Have Installed Software    nodered-demo,0.0.1


Uninstall node-red project via Cumulocity
    # Skip    Missing Uninstall software keyword
    ${operation}=    Cumulocity.Create Operation    fragments={"c8y_SoftwareUpdate":[{"name":"nodered-demo","version":"latest::nodered","url":"","action":"delete"}]}    description=Remove nodered-demo package
    Operation Should Be SUCCESSFUL    ${operation}
    ${mo}=    Cumulocity.Device Should Have Fragments    c8y_SoftwareList
    Log    ${mo}
    Should Not Contain    ${mo}    nodered-demo
