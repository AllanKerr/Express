# Endpoints Configuration Parsing

Endpoints configuration `.yaml` files must be parsed when using the `deploy` and `update` operations. This process involves producing a set of [Kubernetes Ingress](https://kubernetes.io/docs/concepts/services-networking/ingress/) configurations from the file's contents. [Information on the endpoints configuration file specification can be found here.](./endpoints-file.md)

When requests to endpoints with scopes are made, the reverse proxy adds the required scopes as a header and forwards the access token to the authorization service. The reverse proxy is responsible for adding the headers to improve performance by removing the need for a path lookup to determine which scopes are required. However, [Kubernetes Ingress](https://kubernetes.io/docs/concepts/services-networking/ingress/) configurations only allow headers to be added using [annotations](https://github.com/kubernetes/ingress-nginx/blob/master/docs/annotations.md) which apply to all paths in the Ingress configuration. Because of this, an Ingress configuration is required for each distinct set of scopes. While this may sound inefficient, all Ingress configurations are handled by a single [nginx ingress controller](https://github.com/kubernetes/ingress-nginx) meaning there is only a slight increase in overhead during deployment.  

![Image of Yaktocat](./file-parsing.png)

This diagram provides an example of the resulting Ingress configurations from parsing an `Endpoints-Configuration.yaml` file. The set of Ingress configurations is returned as the result of the parsing process.

## 1. Grouping
The first step during parsing is to group the path's based on their set of scopes. This is done by using the hash code for the set of scopes as the key in map. The hash code for the set of scopes is created by hashing each scope and XOR-ing the results. This allows for the same hash code to be produced even if scopes are listed in a different order. This step results in path groups which consist of a set of scopes and a set of paths that require those scopes.

During this step a hash set of all the paths found in the configuration file is found to return an error if any duplicate paths are discovered.

## 2. Ingress Configurations
Each path group must then be converted to an Ingress configuration.

### Labels
Each configuration is given a set of labels. The `app=name` label is given to allow lookup of all Ingress configurations associated with a deployed application container. The `identifier=name:hashcode` label is also added tagging each ingress configuration with the deploy name and the hash code of the scopes set. This allows existing Ingress configurations to be looked up if the endpoints for a deployed application container are modified.
### Annotations
Each configuration is also given a set of annotations. These are required for advanced configuration of the underlying [nginx ingress controller](https://github.com/kubernetes/ingress-nginx).

If the group has one or more scopes then a `configuration-snippet` is added to add the required scopes to the request header and forward the request to the authorization server's introspection endpoint. If introspection succeeds then `proxy_set_header` is used to pass the user id and user scopes that the access token was created for on to the deployed application container. The `configuration-snippet` can be found below.
```
set $scopes group.scopes;
auth_request        /external-auth;
auth_request_set    $auth_cookie $upstream_http_set_cookie;
add_header          Set-Cookie $auth_cookie;
auth_request_set    $authHeader0 $upstream_http_user_id;
proxy_set_header    'User-Id' $authHeader0;
auth_request_set    $authHeader1 $upstream_http_user_scopes;
proxy_set_header    'User-Scopes' $authHeader1;
```
The authorization snippet is only added for groups with one or more required scopes. However, a rewrite `configuration-snippet` is added for all groups. This removes the deploy name from the path before forwarding the request to the deployed application container. This means that if the path is exposed at `https://cluster-ip/example/unprotected` then the application container will receive the request at `/unprotected` instead of `/example/unprotected`. This is accomplished using the following snippet to remove the deploy name.
```
rewrite ^/%v/(.*)$ /$1 break;
```
### Backend
The last component of the Ingress configuration is the backend that the paths are sent to. This requires a name and port to route requests to the deployed `Kubernetes service`.
