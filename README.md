# Tedge node-red plugin

## Pre-requisites

You need to have node-red already installed on your device. If you don't already have it then check out the [Node-RED documentation](https://nodered.org/docs/getting-started/).

## Plugin summary

Install/remove node-red projects on a device using thin-edge.io software management plugin.

**Technical summary**

The following details the technical aspects of the plugin to get an idea what systems it supports.

|||
|--|--|
|**Languages**|`shell` (posix compatible)|
|**CPU Architectures**|`all/noarch`. Not CPU specific|
|**Supported init systems**|`N/A`|
|**Required Dependencies**|-|
|**Optional Dependencies (feature specific)**|-|

### How to do I get it?

The following linux package formats are provided on the releases page and also in the [tedge-community](https://cloudsmith.io/~thinedge/repos/community/packages/) repository:

|Operating System|Repository link|
|--|--|
|Debian/Raspian (deb)|[![Latest version of 'tedge-nodered-plugin' @ Cloudsmith](https://api-prd.cloudsmith.io/v1/badges/version/thinedge/community/deb/tedge-nodered-plugin/latest/a=all;d=any-distro%252Fany-version;t=binary/?render=true&show_latest=true)](https://cloudsmith.io/~thinedge/repos/community/packages/detail/deb/tedge-nodered-plugin/latest/a=all;d=any-distro%252Fany-version;t=binary/)|
|Alpine Linux (apk)|[![Latest version of 'tedge-nodered-plugin' @ Cloudsmith](https://api-prd.cloudsmith.io/v1/badges/version/thinedge/community/alpine/tedge-nodered-plugin/latest/a=noarch;d=alpine%252Fany-version/?render=true&show_latest=true)](https://cloudsmith.io/~thinedge/repos/community/packages/detail/alpine/tedge-nodered-plugin/latest/a=noarch;d=alpine%252Fany-version/)|
|RHEL/CentOS/Fedora (rpm)|[![Latest version of 'tedge-nodered-plugin' @ Cloudsmith](https://api-prd.cloudsmith.io/v1/badges/version/thinedge/community/rpm/tedge-nodered-plugin/latest/a=noarch;d=any-distro%252Fany-version;t=binary/?render=true&show_latest=true)](https://cloudsmith.io/~thinedge/repos/community/packages/detail/rpm/tedge-nodered-plugin/latest/a=noarch;d=any-distro%252Fany-version;t=binary/)|

### What will be deployed to the device?

* The following software management plugins which is called when installing and removing `nodered` projects via Cumulocity IoT
    * `nodered` - Deploy a project using the node-red project structure (e.g. git repository containing a flow)


### Activating a project

Node-red supports deploying multiple projects to a device however only one project can be active at one time. The `tedge-nodered-plugin` supports installing multiple projects, and the active project can be controlled by a special (magic) package name

A project can be activated by using a special package name and version.

|Name|Version|SoftwareType|
|--|--|--|
|`active-project`|`<project_name>`|`nodered`|

The reasoning behind having the magic package was to allow users to deploy projects to devices and decide whether it should be activated or not. Sometimes not automatically activating a project can be useful if you want to download the project to all the devices but wait to activate it at a later point in time.

## Plugin Dependencies

The following packages are required to use the plugin:

* jq
* curl
* node-red
* npm (requirement of node-red)

## Known issues

* Activating a new project for the first time does not seem to work for reasons unknown, however restarting the `nodered` service seems to fix this. The root cause should be investigated as it is probably a sign that an API call is missing or the install sequence is wrong.

## Design

### Proposed workflow of a node-red project

The following workflow details how node-red projects could be deployed to a fleet of devices. It shows the development process through to the deployment to devices.

1. Create git repo for a node-red project (containing a flow)

2. Build/test the flow using a github runner or run on a device

3. Publish the flow to Cumulocity using go-c8y-cli to the software package (with new version)

4. Manually create a bulk operation to deploy the new flow to the fleet of devices

### Package format

#### Option1: Configuration file containing the URL of the project

A single file containing the git URL which can be used to clone the project on the device.

```sh
REPO=https://github.com/reubenmiller/nodered-demo.git
```

#### Option 2: Single tarball (tar.gz) file containing a snapshot of the project

The single tarball which contains the exported project's files.

* flow.json
* settings.js (used to parameterize the flow.json)
* dependencies.txt (text file containing each dependency one per line) this will be used to install npm modules
* version file containing the version of the snapshot on the first line

    ```sh
    1.0.0
    ```

### Packaging a node-red project

Instructions how to create an offline installable node-red project. It will export a git project and create a tarball (`.tar.gz`) that can be used with this plugin.

#### Archive the full directory

```sh
COMMIT=$(git rev-parse --short HEAD);
BRANCH=$(git rev-parse --abbrev-ref HEAD);
VERSION="$BRANCH@$COMMIT";
ARCHIVE="$(basename $(pwd))__$VERSION".tar.gz;
tar cvzf "$ARCHIVE" .
```

#### Archive current branch

1. Create an archive using git (requires a recent version of git which supports the --add-virtual-file flag)

    ```sh
    cd my-project
    COMMIT=$(git rev-parse --short HEAD);
    BRANCH=$(git rev-parse --abbrev-ref HEAD);
    VERSION="$BRANCH@$COMMIT";
    ARCHIVE="$(basename $(pwd))__$VERSION".tar.gz;
    git archive --format=tar.gz --add-virtual-file=version:"$VERSION" -o "$ARCHIVE" HEAD;
    ```

2. Create the software item and upload the version to the Cumulocity software repository

    ```sh
    c8y software create --name nodered-demo --description "Node red flow" --data softwareType=nodered
    ```

    Then add the version

    ```sh
    c8y software versions create --software nodered-demo --version "$VERSION" --file ./nodered-demo__master@32256bb.tar.gz

    # Or using wildcards
    c8y software versions create --software nodered-demo --version "$VERSION" --file ./nodered-demo__*@*.tar.gz
    ```

3. Create an active flow software entry to control which item should be active. There should be one per flow

    ```sh
    c8y software create --name active-flow --description "Active node-red flow" --data softwareType=nodered
    c8y software versions create --software active-flow --version "nodered-demo" --url " "
    ```

### Uploading a project using an external URL

1. Create a configuration file containing the url

    ```sh
    REPO=https://github.com/reubenmiller/nodered-demo
    ```

2. Create a software version item in Cumulocity IoT

    ```sh
    c8y software versions create --software nodered-demo --version "latest" --file ./tests/testdata/nodered-demo.cfg
    ```

## Future ideas

* How to spawn to node-red instances to support applying multiple flows on the same device. Each instance must use a different port / service.
