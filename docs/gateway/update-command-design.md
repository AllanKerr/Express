# Update Command Design

The update command is designed to update the four Kubernetes objects that are created when an application container is deployed. This was designed separately from `deploy` and `teardown` to allow updates without any downtime.

## Update Configurations

To allow for any subset of the parameters to be updated, the design was based around update configurations. An update configuration is a structure that matches the structure of the object configurations used by `deploy`.

However, the update configuration requires pointers for all leaf values rather than the values themselves. This allows each value to convey whether or not the value was updated and, if the value was updated, what it was updated to. A `nil` value indicates that the field was not updated. This design was based on monadic optional values but pointers were used due to golang's lack of support for optionals. `nil` is used in place of an absent optional value.

Two update configurations are offered, the `ContainerUpdate` and the `AutoscalerUpdate`.

```
ContainerUpdate {
	Image *string
}

AutoscalerUpdate {
	MinReplicas *int
	MaxReplicas *int
}
```



`updateDeployment` and `updateAutoscaler` both take the provided parameters and create an update configuration. `updateEndpoints` does not use an update configuration because it uses the [endpoints `.yaml` parser](./endpoints-parsing.md) to produce a new set of Ingress configurations. These configurations may result in the creation of new Ingress objects or updating existing ones depending on whether a distinct set of scopes was added or removed. There is no `updateServices` because the port the service exposes and targets cannot be updated.

## Updaters

All updaters implement the [Updater interface](https://github.com/AllanKerr/Services/blob/master/gateway-controller/gateway-controller/kube/updater.go) which requires `GetModifiers` implementers to return a list of the names of the leaf values that will result in an update and `Update` which takes the update configuration.

### Deployment Updater

The [DeploymentUpdater](https://github.com/AllanKerr/Services/blob/master/gateway-controller/gateway-controller/kube/updater_deployment.go) is designed to rollout new instances if the image is updated. An updated image is the only parameter that will result in a deployment update. The deployment has been designed to accept a `ContainerUpdate` which is an update configuration for the deployment container rather than the Kubernetes deployment as a whole. This design decision was made because only parameters that modify the container are accepted.

The update process occurs inside of a `RetryOnConflict` block. This ensures that if multiple updates occur simultaneously, the one the experiences an update will retry rather than fail. This design is required to eliminate concurrency issues which would result in partial updates and an inconsistent state.

To complete the update, the deployment updater uses the [client-go library](https://github.com/kubernetes/client-go) and the update name to fetch the most recent version of the Kubernetes deployment. The container update is then applied to the fetched configuration. The resulting configuration is then sent using the [client-go library](https://github.com/kubernetes/client-go) to attempt an update. A successful update will result in the [Kubernetes rollout process](https://kubernetes.io/docs/tasks/run-application/rolling-update-replication-controller/) while a failed update will result in a retry.

### Autoscaler Updater

The [AutoscalerUpdater](https://github.com/AllanKerr/Services/blob/master/gateway-controller/gateway-controller/kube/updater_autoscaler.go) is designed to update the minimum and maximum number of instances of the application container. The design uses the `AutoscalerUpdate` which is an update configuration for the Kubernetes HorizontalPodAutoscalerSpec rather than the entire Kubernetes HorizontalPodAutoscaler because the autoscaler metadata should not be updated after being deployed.

The update process also occurs inside of a `RetryOnConflict` block like the Deployment Updater. This uses the [client-go library](https://github.com/kubernetes/client-go) to fetch the most recent version of the Kubernetes Horizontal Pod Autoscaler, apply the `AutoscalerUpdate` configuration and attempt an update. A successful update will alter modify the minimum and/or maximum number of application container instances while a failed update will result in a retry.

### Service Updater

Because none of the supported parameters modify the Kubernetes service, there is no service updater.

### Ingress Updater

The [IngressUpdater](https://github.com/AllanKerr/Services/blob/master/gateway-controller/gateway-controller/kube/updater_ingress.go) is designed to accept a set of Ingress configurations rather than an update configuration. This design choice was made because the new set of paths may result in Ingress configurations being created, updated, or deleted.

#### 1. Port

Because the port is not contained in the set of Ingress configurations, the port must be fetched from the deployed Kubernetes service using the [client-go library](https://github.com/kubernetes/client-go).

#### 2. Existing Ingresses

When updating the exposed endpoints, Ingress configurations may be created, updated, or deleted. This means that the current set of Ingress configurations must be fetched to determine which ones should be updated or deleted. This is done using the [client-go library](https://github.com/kubernetes/client-go) and the `app=name` label set during `deploy` to fetch all current Ingress configurations.

#### 3. New Ingresses

Once the set of existing Ingresses has been fetched, the new set of Ingress configurations can be applied. For each new ingress, the `identifier` label is used to find the existing ingress definition inside the existing ingresses list.

If an existing Ingress configuration is found with a matching identifier, then the existing Ingress is removed from the existing Ingresses list and an update is performed inside a `RetryOnConflict` block using the [client-go library](https://github.com/kubernetes/client-go).

If an existing Ingress configuration isn't found then the Ingress configuration has a new set of scopes that didn't previously exist. This allows for the new Ingress configuration to be created directly as it would during `deploy` using the [client-go library](https://github.com/kubernetes/client-go).

After processing all of the new ingresses, all updates have been performed. However, this does not account for endpoints that have been removed. To remove endpoints that no longer exist, the existing Ingresses list is used because it only contains Ingress configurations with identifiers that were not found in the set of new Ingress configurations. To remove these endpoints, the remaining Ingress configurations are iterated through and removed using the [client-go library](https://github.com/kubernetes/client-go).
