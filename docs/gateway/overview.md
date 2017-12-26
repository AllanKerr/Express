# Gateway Architecture

The gateway is the entry point into the system for all public traffic. It is responsible for load balancing, routing and protecting endpoints. This diagram shows the system components that the gateway interacts with.

*TODO: Add architecture diagram*

## Gateway Controller

The gateway controller provides the command line interface for developers to deploy, update and tear down their application containers. This controller is responsible for automatically reconfiguring the gateway giving developers the ability to rapidly deploy and update applications without the need to create or maintain configuration files.

### Deployment Components
When developers deploy application containers using the gateway controller, four components are automatically created and configured:

*TODO: Add deployment components diagram*

1. *Deployment.* A [Kubernetes Deployment](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/) is created to manage the application container. This allows multiple instances of the container to seamlessly run across multiple machines, deployed application containers to be updated to different [Docker Images](https://docs.docker.com/get-started/part2/) with no downtime and for the deployed application to be automatically restarted in the event of an unexpected failure.
2. *Autoscaler.* A [Kubernetes Horizontal Pod Autoscaler](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/) is created to automatically scale the deployment based on CPU utilization. The autoscaler is responsible for monitoring the application instances created by the deployment and creating new instances if CPU utilization exceeds 80%.
3. *Service.* A [Kubernetes Service](https://kubernetes.io/docs/concepts/services-networking/service/) is created to reference and provide access to the instances managed by the deployment. This exposes the IP addresses of the instances to allow for load balancing and routing.
4. *Gateway.* A set of [Kubernetes Ingresses](https://kubernetes.io/docs/concepts/services-networking/ingress/) are created to expose the service to the public. An ingress configuration is created for each set of endpoint access scopes. A default configuration that routes directly to the service is created for unprotected endpoints while configurations that route to the authorization service's introspection endpoint are created for protected endpoints.

### Deployment Operations

The automatic configuration of these four components is handled through four operations that developers can perform on the gateway controller.

1. *Deploy.* Deploys a new application container in the form of a [Docker Image](https://docs.docker.com/get-started/part2/). This is used by developers to instantly deploy their application containers and results in the automatic configuration and creation of all four deployment components mentioned above.
2. *List.* Lists the application containers deployed using the `deploy` operation. This can be used by developers to view the applications that have been deployed to the system.
3. *Update.* Updates a specified application container found using the `list` operation. This allows for developers to deploy new versions of their application containers with zero downtime.
4. *Teardown.* Teardown a deployed application container found using the `list` operation. This is used by developers to remove application containers that were previously deployed.

## Gateway
