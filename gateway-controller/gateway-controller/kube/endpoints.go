package kube

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	"gopkg.in/fatih/set.v0"
	"errors"
	"k8s.io/apimachinery/pkg/util/intstr"
	"strings"
	"fmt"
)

type Endpoint struct {
	Path string     `yaml:"path"`
	Scopes []interface{} `yaml:"scopes,omitempty"`
}

type EndpointsConfig struct {
	Endpoints []Endpoint `yaml:"endpoints"`
}

type endpointGroup struct {
	scopes *set.Set
	paths []string
}

func (group *endpointGroup) append(path string) {
	group.paths = append(group.paths, path)
}


func (group *endpointGroup) getScopes() string {

	var scopes string
	group.scopes.Each(func(item interface{}) bool {
		if scopes == "" {
			scopes += item.(string)
		} else {
			scopes += " " + item.(string)
		}
		return true
	})
	return scopes
}

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

const RewriteSnippet = `
	rewrite ^/%v/(.*)$ /$1 break;
`

func (group *endpointGroup) getHashCode() int {

	var hashCode int
	group.scopes.Each(func(item interface{}) bool {
		hashCode ^= hashString(item.(string))
		return true
	})
	return hashCode
}

func (group *endpointGroup) getAnnotations(name string) map[string]string {

	var snippet string
	if group.scopes != nil {
		snippet += fmt.Sprintf(AuthSnippet, group.getScopes())
	}
	snippet += fmt.Sprintf(RewriteSnippet, name)
	return map[string]string{
		"ingress.kubernetes.io/configuration-snippet": snippet,
	}
}

func (group *endpointGroup) GetIngress(name string, port int32) *extensionsv1beta1.Ingress {

	labels := map[string]string{
		"app": name,
		"identifier" : string(group.getHashCode()),
	}

	backend := extensionsv1beta1.IngressBackend{
		ServiceName: name,
		ServicePort: intstr.FromInt(int(port)),
	}
	var paths []extensionsv1beta1.HTTPIngressPath
	for _, path := range group.paths {

		if !strings.HasPrefix(path, "/") {
			path += "/"
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
	return &extensionsv1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: prefix,
			Labels: labels,
			Annotations: group.getAnnotations(name),
		},
		Spec: extensionsv1beta1.IngressSpec{
			Rules: []extensionsv1beta1.IngressRule{
				{
					IngressRuleValue: extensionsv1beta1.IngressRuleValue{
						HTTP: &extensionsv1beta1.HTTPIngressRuleValue {
							Paths: paths,
						},
					},
				},
			},
		},
	}
}

func ParseConfig(name string, port int32, config *EndpointsConfig) ([]*extensionsv1beta1.Ingress, error) {

	paths := set.New()
	defaultGroup := &endpointGroup{}
	var groups []*endpointGroup

	for _, endpoint := range config.Endpoints {

		if paths.Has(endpoint.Path) {
			return nil, errors.New("duplicate path in config: " + endpoint.Path)
		}

		if len(endpoint.Scopes) == 0 {
			defaultGroup.append(endpoint.Path)
		} else {
			hasGroup := false
			for _, group := range groups {

				if group.scopes.Has(endpoint.Scopes...) {
					group.append(endpoint.Path)
					hasGroup = true
					break
				}
			}
			if !hasGroup {
				group := &endpointGroup{
					scopes: set.New(endpoint.Scopes...),
					paths: []string{endpoint.Path},
				}
				groups = append(groups, group)
			}
		}
		paths.Add(endpoint.Path)
	}

	var ingresses []*extensionsv1beta1.Ingress
	ingresses = append(ingresses, defaultGroup.GetIngress(name, port))
	for _, grp := range groups {
		ingresses = append(ingresses, grp.GetIngress(name, port))
	}
	return ingresses, nil
}