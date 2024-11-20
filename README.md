# tedge-nodered-plugin

## Pre-requisites

Naturally, node-red must be installed in order to use this plugin as the plugin uses the node-red REST api to managed node-red flows. We recommend installing node-red in a container using the [tedge-container-plugin-ng](https://github.com/thin-edge/tedge-container-plugin/tree/next) thin-edge.io software management plugin.

node-red 

node-red supports two different modes, one is the classic mode when flows are just simple json files (e.g. `flows.json`), and the other is the project mode where a flows.json is deployed via a git repository. The former (simple json files) is more flexible and widely used, so it is the recommended way to deploy flows to a device.

To help with the installation the following docker-compose.yaml files can be used to deploy a node-red container via thin-edge.io.

* [node-red (flows/classic mode))](tests/testdata/docker-compose.nodered-flows.yaml) (recommended)
* [node-red in project mode)](tests/testdata/docker-compose.nodered-project.yaml)

Be sure to check out the [Node-RED documentation](https://nodered.org/docs/getting-started/) for more details on how to configure the node-red container.

## Plugin summary

Install/remove node-red flows or projects on a device using the thin-edge.io software management plugin API.

**Technical summary**

The following details the technical aspects of the plugin to get an idea what systems it supports.

|||
|--|--|
|**Languages**|`golang`|
|**CPU Architectures**|`armv6 (armhf)`, `armv7 (armhf)`, `arm64 (aarch64)`, `amd64 (x86_64)`|
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
    * `nodered-project` - Deploy a project using the node-red project structure (e.g. git repository containing a flow)
    * `nodered-flows` - Deploy a node-red flow (e.g. `flows.json`)

## Plugin Dependencies

The following packages are required to use the plugin:

* node-red (we recommend deploying it as a container)


### Deploying format

#### nodered-flows

A node-red flows file, is the classic node-red json format which you get when you export the node-red project from the node-red UI.

Example flows:

* [flows.json](https://github.com/reubenmiller/nodered-demo-next/blob/main/flows.json)

You can use [go-c8y-cli](https://goc8ycli.netlify.app/) to create the Cumulocity IoT software repository items for your flow:

```sh
# Create a new software item
c8y software create --name myflow --softwareType nodered-flows

# For each version, upload a new flows.json file
wget -O - https://github.com/reubenmiller/nodered-demo-next/blob/main/flows.json > flows.json
c8y software versions create --software myflow --version 1.0.0 --file ./flows.json
```


#### nodered-project

A node-red project can be deployed to a device via the software management feature, where the software artifact is a simple json format which the `.repo` property which indicates the Git repository of the node-red project which should be deployed to the device.

Below is an example of such as deployment artifact.

```json
{
    "repo": "https://github.com/reubenmiller/nodered-demo-next"
}
```

You can use [go-c8y-cli](https://goc8ycli.netlify.app/) to create the Cumulocity IoT software repository items for your flow:

```sh
# Create a new software item
c8y software create --name my-nodered-project --softwareType nodered-project

# For each version, upload a new flows.json file
echo '{"repo": "https://github.com/reubenmiller/nodered-demo-next"}' > my-nodered-project.json
c8y software versions create --software my-nodered-project --version 1.0.0 --file ./my-nodered-project.json
```
