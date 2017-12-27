# Endpoints Configuration Specification

When developers deploy application containers, there is a set of endpoints that they want to expose to the public. The design choice was made to use a configuration file over a client library to allow developers to expose endpoints for application containers written in any language. This choice was also made to avoid introducing coupling between the Services project and the application containers deployed using it.

# File Specification
To minimize the time spent creating configuration files, a simple `.yaml` specification was designed to define exposed endpoints. This consists of the `endpoints` key and an array of paths as the value. When endpoints are exposed, they are always prefixed with the name provided to the `deploy` operation to avoid path conflicts.

## Examples
### Basic
Here is a basic specification that exposes the two endpoints `/examplepath1` and `/examplepath2`.
```
endpoints:
  - path: /examplepath1
  - path: /examplepath2
```
If this specification is used when deploying an application container named `example` then the final public paths will be:
1. `https://cluster-ip/example/examplepath1`
2. `https://cluster-ip/example/examplepath2`

### Advanced
Here is a a simple specification that exposes the three endpoints `/unprotected`, `/protected/user` and `/protected/admin`.
```
endpoints:
  - path: /unprotected
  - path: /protected/user
    scopes: [user]
  - path: /protected/admin
    scopes: [user, admin]
```
If this specification is used when deploying an application container named `example` then the final public paths will be:
1. `https://cluster-ip/example/unprotected`
2. `https://cluster-ip/example/protected/user`
3. `https://cluster-ip/example/protected/admin`

The `/unprotected` endpoint can be accessed normally because it is unprotected. However, accessing `/protected/user` or `/protected/admin` will result in 401 unauthorized. These two endpoints require a valid OAuth2 Access token with all the matching scopes. Scopes provide a simple method for developers to protect their endpoints. Direct matching was used for scopes rather than a hierarchical or regular expression approach for performance reasons.

The first endpoint has the scope `user` meaning the access token must also have the user scope for introspection to succeed.
The second endpoint has the scopes `user` and `admin`. The access token must possess both scopes for introspection to succeed. A token with only the `user` scope or only the `admin` will still receive 401 unauthorized.  
