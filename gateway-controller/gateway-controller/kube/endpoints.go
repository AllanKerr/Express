package kube

import (
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	"gopkg.in/fatih/set.v0"
	"errors"
	"k8s.io/apimachinery/pkg/util/intstr"
	"strings"
	"fmt"
	"strconv"
	"gopkg.in/yaml.v2"
)

// The structure for an endpoint for the endpoints config file
type Endpoint struct {
	Path string     `yaml:"path"`
	Scopes []string `yaml:"scopes,omitempty"`
}

// The base structure for the endpoints config file
type EndpointsConfig struct {
	Endpoints []Endpoint `yaml:"endpoints"`
}

// A group of endpoints or paths with the same set of scopes regardless of order
// Used to group paths that share scopes during parsing
type endpointGroup struct {
	scopes *set.Set
	paths []string
}

// Add a path to the endpoint group
func (group *endpointGroup) append(path string) {
	group.paths = append(group.paths, path)
}

// Get the scopes string in an undefined order
// Appends all scopes separated by a single space between each
func (group *endpointGroup) getScopes() string {

	var scopes string
	group.scopes.Each(func(item interface{}) bool {
		if scopes == "" {
			// no space before first scope
			scopes += item.(string)
		} else {
			scopes += " " + item.(string)
		}
		return true
	})
	return scopes
}

// The configuration snippet for scoped endpoints to redirect to the authorization service
// Checks both the OAuth2 access token and the set of required scopes
const AuthSnippet = `
	# this location requires authentication
	set $scopes "%v";
	auth_request        /external-auth;
	auth_request_set    $auth_cookie $upstream_http_set_cookie;
	add_header          Set-Cookie $auth_cookie;
	auth_request_set $authHeader0 $upstream_http_user_id;
	proxy_set_header 'User-Id' $authHeader0;
	auth_request_set $authHeader1 $upstream_http_user_scopes;
	proxy_set_header 'User-Scopes' $authHeader1;
`

// The configuration snippet used by all configurations to remove the deployment name
// from the path before proxy passing the request to the internal deployment
const RewriteSnippet = `
	rewrite ^/%v/(.*)$ /$1 break;
`

// Gets a hash code for an endpoint group based on its set of scopes
// The hash code will be the same even if the order of the scopes is different
func (group *endpointGroup) getHashCode() int {

	if group.scopes == nil {
		return HashString("")
	}
	var hashCode int
	group.scopes.Each(func(item interface{}) bool {
		hashCode ^= HashString(item.(string))
		return true
	})
	return hashCode
}

// Gets the annotations that will be added to the Ingress configuration
func (group *endpointGroup) getAnnotations(name string) map[string]string {

	var snippet string
	// Add auth snippet if the endpoint is scoped meaning it needs protection
	if group.scopes != nil && !group.scopes.IsEmpty() {
		snippet += fmt.Sprintf(AuthSnippet, group.getScopes())
	}
	snippet += fmt.Sprintf(RewriteSnippet, name)
	return map[string]string{
		"nginx.ingress.kubernetes.io/configuration-snippet": snippet,
		"kubernetes.io/ingress.class": "nginx",
	}
}

// Get an Ingress configuration from an endpoint group with name
// being the name of the deployed service the endpoints are routed to
// and port being the port the deployed service listens on
func (group *endpointGroup) GetIngress(name string, port int32) *extensionsv1beta1.Ingress {

	backend := extensionsv1beta1.IngressBackend{
		ServiceName: name,
		ServicePort: intstr.FromInt(int(port)),
	}
	var paths []extensionsv1beta1.HTTPIngressPath
	for _, path := range group.paths {

		if !strings.HasPrefix(path, "/") {
			path = "/" + path
		}
		paths = append(paths, extensionsv1beta1.HTTPIngressPath{
			Backend: backend,
			Path: "/" + name + path,
		})
	}
	prefix := name
	if !strings.HasSuffix(name, "-") {
		prefix += "-"
	}

	ingress := DefaultIngressConfig()

	ingress.ObjectMeta.GenerateName = prefix
	ingress.Annotations = group.getAnnotations(name)
	ingress.Labels = map[string]string{
		"app": name,
		"identifier" : strconv.Itoa(group.getHashCode()),
	}

	ingress.Spec.Rules[0].HTTP.Paths = paths
	return ingress
}

// unmarshalls a endpoints config file into an EndpointsConfig struct
func unmarshallConfig(contents []byte) (*EndpointsConfig, error) {
	var endpoints EndpointsConfig
	if err := yaml.Unmarshal(contents, &endpoints); err != nil {
		return nil, err
	}
	return &endpoints, nil
}

// Parses an endpoints configuration file into a set of Ingress configurations
// Each Ingress configuration routes the path to the deployed service with
// name and port
//
// The file specification can be found here:
// https://github.com/AllanKerr/Express/blob/master/docs/gateway/endpoints-file.md
func ParseConfig(name string, port int32, contents []byte) ([]*extensionsv1beta1.Ingress, error) {

	config, err := unmarshallConfig(contents)
	if err != nil {
		return nil, err
	}

	paths := set.New()
	defaultGroup := &endpointGroup{}
	var groups []*endpointGroup

	for _, endpoint := range config.Endpoints {

		// error if a duplicate path is found
		if paths.Has(endpoint.Path) {
			return nil, errors.New("duplicate path in config: " + endpoint.Path)
		}

		// add to the default group if it has no scopes meaning it doesn't need protection
		if len(endpoint.Scopes) == 0 {
			defaultGroup.append(endpoint.Path)
		} else {
			hasGroup := false

			scopes := stringSliceToInterfaceSlice(endpoint.Scopes)

			// check for an existing group to add the endpoint to
			// that shares the same set of scopes
			for _, group := range groups {
				if group.scopes.Has(scopes...) {
					group.append(endpoint.Path)
					hasGroup = true
					break
				}
			}
			// create a new group for the endpoint if it has a new set of scopes
			if !hasGroup {
				group := &endpointGroup{
					scopes: set.New(scopes...),
					paths: []string{endpoint.Path},
				}
				groups = append(groups, group)
			}
		}
		paths.Add(endpoint.Path)
	}

	// create an Ingress configuration for each endpoint group
	var ingresses []*extensionsv1beta1.Ingress
	if len(defaultGroup.paths) > 0 {
		ingresses = append(ingresses, defaultGroup.GetIngress(name, port))
	}
	for _, grp := range groups {
		ingresses = append(ingresses, grp.GetIngress(name, port))
	}
	return ingresses, nil
}

// convert a slice of strings to a slice of interfaces
func stringSliceToInterfaceSlice(strs []string) []interface{} {
	s := make([]interface{}, len(strs))
	for i, v := range strs {
		s[i] = v
	}
	return s
}
