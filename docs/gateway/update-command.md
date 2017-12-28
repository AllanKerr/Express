# Update Command

The update command allows developers to update previously deployed application containers. This allows for the Docker image to be changed, the set of exposed endpoints to be reconfigured, and the minimum and maximum number of instance replications to be changed. This could be achived using `teardown` and `deploy`; however, `update` allows the application container to be updated without any downtime.

## Command

The update command is made available to developers through the `services` executable.

```
services update <name> [flags]

Flags:
      --endpoint-config string   The configuration file for accessing the deployed API.
      --image string             The new Docker image to roll out.
      --max int32                The minimum number of instances.
      --min int32                The minimum number of instances.
```

## Parameters

The `update` command allows all parameters defined in the `deploy` command to be modified except for port. The only way to update the port of a deployed application container is to use `teardown` followed by `deploy` which will result in some downtime.

1. ***Name.*** The name of the application container to be updated must be specified. This must match the name passed to the `deploy` command and be present when using the `list` command. Nothin is changed if the name does not match an existing deployment.
2. ***Endpoints Configuration.*** The endpoints configuration `.yaml` file may be updated to expose, hide or change the required scopes of endpoints.
3. ***Image.*** The [Docker Image](https://docs.docker.com/get-started/part2/) of the application may be updated. This may be a new version of the same image or a different image but must listen on the original port. If the image is changed, the instances of the new application container are brought up before the old version is taken down to ensure zero downtime occurs.
4. ***Minimum Replicas.*** The minimum number of instances of the application container to deploy. This is enforced by the autoscaler to ensure that the number of available instances never drops below this number regardless of CPU utilization.
5. ***Maximum Replicas.*** The maximum number of instances of the application container to deploy. This is enforced by the autoscaler to ensure that the number of available instances never exceeds this number regardless of CPU utilization.

## Access

Once an application container has been updated, it may take up some time before the changes take effect. The changes will take the longest to take effect if the image is changed because it must be pulled from the remote Docker repository.

Although the Services project does not support checking when the status of the rollout, this can be monitored using the Kubernetes `kubectl rollout status deployment <name>` command.

## Design

[The design of the update command can be found here.](./update-command-design.md)
