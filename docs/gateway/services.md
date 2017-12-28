# Services

The `services` executable is offered to allow developers to deploy, list, update and teardown application containers using the Services project. This is offered as part of the `gateway-controller`.

## Building
The `services` executable can be built manually from the [`gateway-controller`](https://github.com/AllanKerr/Services/tree/master/gateway-controller) directory using the `make` command. Manual building requires `golang` to be installed on the system.

## Commands
The `services` executable offers four commands for automatically configuring the gateway to run application containers.

### 1. Deploy
Deploys a new application container in the form of a [Docker Image](https://docs.docker.com/get-started/part2/). This is used by developers to instantly deploy their application containers and results in the automatic configuration and creation of all four deployment components mentioned above.

**[Deploy command documentation can be found here.](./deploy-command.md)**

### 2. List
Lists the application containers deployed using the `deploy` operation. This can be used by developers to view the applications that have been deployed to the system.

**[List command documentation can be found here.](./list-command.md)**

### 3. Update
Updates a specified application container found using the `list` operation. This allows for developers to deploy new versions of their application containers with zero downtime.

**[Update command documentation can be found here.](./update-command.md)**

### 4. Teardown
Teardown a deployed application container found using the `list` operation. This is used by developers to remove application containers that were previously deployed.

**[Teardown command documentation can be found here.](./teardown-command.md)**
