# Node-RED notes

## node-red API

### Check active projects

```sh
curl http://127.0.0.1:1880/projects
```

```json
{"projects":["my-flow","nodered-demo","second-project"],"active":"second-project"}
```

### Activate a project

```sh
curl -H "Content-Type: application/json" -XPUT http://localhost:1880/projects/nodered-demo -d '{"active":true,"clearContext":true}'
```

### Create a new project (need to verify this)

```sh
curl -H "Content-Type: application/json" -XPOST http://127.0.0.1:1880/projects -d '{"action":"create","name":"nodered-demo","summary":"","files":{"flow":"flow.json","credentials":"flow_cred.json"},"migrateFiles":true,"credentialSecret":false}'
```

## Misc.

Node-red's project feature can also be activated by setting an environment variable. It is not used in this project as the demo uses a multi-process container with node-red running under systemd.

```sh
NODE_RED_ENABLE_PROJECTS=true
```

## References

* https://qbee.io/docs/deploy-node-red-flow-github.html
* https://qbee.io/docs/qbee-node-red-deployment.html
