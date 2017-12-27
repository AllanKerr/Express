# Deploy Operation

The deploy operation is designed to allow developers to deploy application containers with minimal need for configuration. This allows developers to rapidly deploy new applications without any of the overhead associated with directly managing Kubernetes objects.

## Parameters

The `deploy` operation was designed to accept the minimum number of parameters required to run, scale, and expose application containers. Because the `deploy` operation creates underlying Kubernetes objects, advanced users can modify their deployed application containers through [kubectl](https://kubernetes.io/docs/reference/kubectl/overview/) for more advanced configuration.

1. *Name.* All deployed applications must be given a name by the developer. This is used to identify the Kubernetes objects associated with the deployed application. All Kubernetes objects are given the label `app=name` to allow identification. Because of this, all deployed applications must have unique names.
2. *Image.* All deployed applications must be given a [Docker Image](https://docs.docker.com/get-started/part2/) which is executed to run the application. The current system is limited to images on public Docker repositories because images are pulled when the application is first deployed. Support for private Docker repositories will be added in the future.
3. *Port.* The port of the application that the Docker container uses may be specified. All traffic to the deployed application will be sent through this port. If not port is provided, it will default to port 80.
4. *Minimum Replicas.* The minimum number of instances of the application container to deploy. This is enforced by the autoscaler to ensure that the number of available instances never drops below this number regardless of CPU utilization.
5. *Maximum Replicas.* The maximum number of instances of the application container to deploy. This is enforced by the autoscaler to ensure that the number of available instances never exceeds this number regardless of CPU utilization.
6. *Endpoints Configuration.* To expose the application to the public through the gateway, an endpoints configuration `.yaml` file must be provided. This allows the developer to specify a set of paths to be exposed along with the a set of scopes on a per-path basis that the Oauth2 access token must possess to access the endpoint. [The design of the endpoints configuration specification is detailed here.](./endpoints.md)

## Design

The design of the deploy operation is based on a transaction model where a transaction is used for each of the four Kubernetes objects. Successful transactions are tracked by the deployment handler. This allows for the handler to rollback all transactions if any of them fail.

## Deploy Handler

The [deploy handler](https://github.com/AllanKerr/Services/blob/master/gateway-controller/gateway-controller/cmd/handlers/handler_deploy.go) is responsible for managing managing the deploy operation to ensure that the entire deploy is successful or none of it is performed.  

The deploy handler is given a reference to the [Kube Client](https://github.com/AllanKerr/Services/blob/master/gateway-controller/gateway-controller/kube/client.go) when it is created. This client uses the local configuration for the running Kubernetes cluster to allow the deploy handler to communicate with the cluster through the [client-go library](https://github.com/kubernetes/client-go).


### 1. Create Service

Once the deploy handler has been created, `createService` is called. This creates the [Kubernetes Service](https://kubernetes.io/docs/concepts/services-networking/service/) using the default service configuration and the provided name and port. The default service configuration specifies a [node port service](https://kubernetes.io/docs/concepts/services-networking/service/#type-nodeport) with a single TCP port.

```
Service {
	ObjectMeta: {}
	Spec: {
		Ports[]: [{Protocol: TCP}]
		Type: NodePort
	}
}
```
Inside `createService`, the default configuration is copied and the deploy name, port and selector are added. The selector is used to specify the set of instances that the service references in the form of `app=name`. An additional `group` label is added when services are created to differentiate between internal services and services that belong to deployed applications. Deployed applications can be found by looking at all services with the label `group=services`.

### 2. Create Deployment

Once the service has been created, `createDeployment` is called. This creates the [Kubernetes Deployment](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/) which will run the application instances referenced by the service. Deployments are also created from a default configuration.  
The default configuration contains the specification for the [replication controller](https://kubernetes.io/docs/concepts/workloads/controllers/replicationcontroller/) that is responsible for managing instances of the application container and the [pod specification](https://kubernetes.io/docs/concepts/workloads/pods/pod/) which contains the application container. The container specification exposes a single port for the service to connect to.  
The container is also given a pre-stop command to sleep for 15 seconds. This delay is required for updating deployed application containers with zero downtime because of the update interval for the reverse proxy before the set of IP addresses for the new updated application container are detected. [A detailed discussion can be found here](https://github.com/kubernetes/ingress-nginx/issues/322). An external PreStop command was used rather than a SIGTERM handler with a delay to minimize coupling with the application container.

```
Deployment {
	ObjectMeta: {}
	Spec: {
		Selector: {}
		Template: {
			ObjectMeta: ObjectMeta{}
			Spec: {
				Containers[]: [{
						Ports[]: [{Name: "http", Protocol: TCP}]
						Lifecycle: {
							PreStop: {
								Exec: {
									Command: {"sleep", "15"},
								}
							}
						}
					}]
				}
			}
		}
	}
}
```
Inside `createService`, the default configuration is modified to give both the replication controller and pod specification the `app=name` label for identification. The `app=name` label is vital because it matches with the service's selector to allow the service to find all instances. The container specification is also updated to use the Docker image specified by the developer and to expose the same port that the service is expecting. The Kubernetes deployment is created with a replication factor of one because the number of replicas is delegated to the autoscaler.

### 3. Create Autoscaler
After the deployment has been created, an autoscaler is created by calling `createAutoscaler` to add or remove instances of the containerized application based on CPU utilization. This results in the creation of a [Kubernetes Horizontal Pod Autoscaler](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/).   
A basic default configuration is used only specifying that the autoscaler will be used to scale a **Kubernetes** deployment.
```
HorizontalPodAutoscaler {
	ObjectMeta: {}
	Spec: {
		ScaleTargetRef: {
			Kind: "Deployment"
	    }
    }
}
```
Inside `createAutoscaler`, the name of the deployment and the `app=name` name label are added to allow the autoscaler to locate the **Kubernetes** deployment it was created to manage. The minimum and maximum number of replicas are also provided. **If no user input was given then these are both set to 1**.

### 4. Create Endpoints
The last step is to create the endpoints that are exposed to the public to allow access to the deployed application container. This process uses the endpoints .yaml parser to produce the set of [Kubernetes Ingress](https://kubernetes.io/docs/concepts/services-networking/ingress/) configurations [as detailed here](./endpoints.md).
