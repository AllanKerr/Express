# List Command Design

The design of the `deploy` command is explicitly designed to support the `list` command. Kubernetes labels are used to allow the set of deployed application containers to be listed using the [client-go library](https://github.com/kubernetes/client-go).

## Kube Client

The [`Kube Client`](https://github.com/AllanKerr/Services/blob/master/gateway-controller/gateway-controller/kube/client.go) acts as a wrapper around the `client-go` Kubernetes client. `ListServices` was added to give the client the ability to lookup all deployed application containers.

`ListServices` fetches a list of all Kubernetes Services with the label `group=services`. This label was added in the `deploy` command to support this. Each deployed container consists of a Kubernetes deployment, service, horizontal pod autoscaler, and ingress configurations. The decision to give the Kubernetes service the `group=services` label was an arbitrary decision. The Kubernetes deployment or horizontal pod autoscaler could have been used because the design only requires one of each per deployed application container.
