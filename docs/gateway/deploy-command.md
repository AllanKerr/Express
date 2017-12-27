# Deploy Command

The deploy operation is designed to allow developers to deploy application containers with minimal need for configuration. This allows developers to rapidly deploy new applications without any of the overhead associated with directly managing Kubernetes objects.

## Command

The deploy command is made available to developers through the `services` executable.

```
services deploy <name> <image> [flags]

Flags:
    --endpoint-config string   The configuration file for protecting the deployed API.
    --max int32                The minimum number of instances. (default 1)
    --min int32                The minimum number of instances. (default 1)
    --port int32               The port exposed by the Docker image. (default 80)
```

## Parameters

The `deploy` operation was designed to accept the minimum number of parameters required to run, scale, and expose application containers. Because the `deploy` operation creates underlying Kubernetes objects, advanced users can modify their deployed application containers through [kubectl](https://kubernetes.io/docs/reference/kubectl/overview/) for more advanced configuration.

1. *Name.* All deployed applications must be given a name by the developer. This is used to identify the Kubernetes objects associated with the deployed application. All Kubernetes objects are given the label `app=name` to allow identification. Because of this, all deployed applications must have unique names.
2. *Image.* All deployed applications must be given a [Docker Image](https://docs.docker.com/get-started/part2/) which is executed to run the application. The current system is limited to images on public Docker repositories because images are pulled when the application is first deployed. Support for private Docker repositories will be added in the future.
3. *Port.* The port of the application that the Docker container uses may be specified. All traffic to the deployed application will be sent through this port. If not port is provided, it will default to port 80.
4. *Minimum Replicas.* The minimum number of instances of the application container to deploy. This is enforced by the autoscaler to ensure that the number of available instances never drops below this number regardless of CPU utilization.
5. *Maximum Replicas.* The maximum number of instances of the application container to deploy. This is enforced by the autoscaler to ensure that the number of available instances never exceeds this number regardless of CPU utilization.
6. *Endpoints Configuration.* To expose the application to the public through the gateway, an endpoints configuration `.yaml` file must be provided. This allows the developer to specify a set of paths to be exposed along with the a set of scopes on a per-path basis that the Oauth2 access token must possess to access the endpoint. [The design of the endpoints configuration specification is detailed here.](./endpoints-file.md)

## Access

Once an application container has been deployed, it can be accessed on any of the endpoints specified in the endpoints configuration file. If the path `/testpath` is declared in the endpoints configuration specification then the application container can be accessed from `https://cluster-ip/name/testpath`.

When a request is passed to the application container, the path is rewritten to remove the deploy name. A request to `https://cluster-ip/name/testpath` will be received on `/testpath` on the application container.
