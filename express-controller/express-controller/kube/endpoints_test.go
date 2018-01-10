package kube

import (
	"testing"
	"gopkg.in/fatih/set.v0"
	"strings"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	"strconv"
)

// Test printing the combined scopes for an endpoint group
func TestGroupGetScopes(t *testing.T) {

	grp := endpointGroup{}
	grp.scopes = set.New()

	if grp.getScopes() != "" {
		t.Errorf("Unexpected empty scopes string: %v", grp.getScopes())
	}

	grp.scopes.Add("first")

	if grp.getScopes() != "first" {
		t.Errorf("Unexpected single scopes string: %v", grp.getScopes())
	}

	grp.scopes.Add("second")

	if grp.getScopes() != "first second" {
		t.Errorf("Unexpected multiple scopes string: %v", grp.getScopes())
	}
}

// Test appending a path to an endpoint group
func TestGroupAppend(t *testing.T) {

	grp := endpointGroup{}

	grp.append("")

	if len(grp.paths) != 1 || grp.paths[0] != "" {
		t.Errorf("Error appending path \"%v\"", "")
	}

	grp.append("second")

	if len(grp.paths) != 2 || grp.paths[0] != "" || grp.paths[1] != "second" {
		t.Errorf("Error appending path \"%v\"", "second")
	}
}

// Test the hash code generation for an endpoints group
func TestGroupHashcode(t *testing.T) {

	grp := endpointGroup{}

	// Test nil scopes
	if grp.getHashCode() != HashString("") {
		t.Error("Error generating hashcode for group with no scopes")
	}

	// Test different orders
	grp1 := endpointGroup{}
	grp1.scopes = set.New()
	grp1.scopes.Add("a")
	grp1.scopes.Add("b")
	grp1.scopes.Add("c")

	grp2 := endpointGroup{}
	grp2.scopes = set.New()
	grp2.scopes.Add("c")
	grp2.scopes.Add("b")
	grp2.scopes.Add("a")

	if grp1.getHashCode() != grp2.getHashCode() {
		t.Error("Error, hash codes did not match for groups with different scope orders")
	}

	// Test duplicate scopes
	grp1.scopes.Add("a")
	grp2.scopes.Add("b")

	if grp1.getHashCode() != grp2.getHashCode() {
		t.Error("Error, hash codes did not match for groups with different scope orders")
	}

	// Test different scopes
	grp1.scopes.Add("diff")
	grp2.scopes.Add("erent")

	if grp1.getHashCode() == grp2.getHashCode() {
		t.Error("Error, expected different hash codes for different scope sets")
	}
}

func TestGroupAnnotations(t *testing.T) {

	grp := endpointGroup{}
	var snippet string
	var ok bool

	// test that scope-less groups only have the rewrite snippet
	annotations := grp.getAnnotations("service")
	if snippet, ok = annotations["ingress.kubernetes.io/configuration-snippet"]; !ok {
		t.Error("Error, annotations missing configuration-snippet")
	}
	if strings.Contains(snippet, "auth_request") {
		t.Error("Error, scope-less endpoints group contains auth snippet")
	}
	if !strings.Contains(snippet, "rewrite") {
		t.Error("Error, expected rewrite snippet")
	}

	grp.scopes = set.New()

	// test that scope-less non-nil scope groups only have the rewrite snippet
	annotations = grp.getAnnotations("service")
	if snippet, ok = annotations["ingress.kubernetes.io/configuration-snippet"]; !ok {
		t.Error("Error, annotations missing configuration-snippet")
	}
	if strings.Contains(snippet, "auth_request") {
		t.Error("Error, scope-less endpoints group contains auth snippet")
	}
	if !strings.Contains(snippet, "rewrite") {
		t.Error("Error, expected rewrite snippet")
	}

	grp.scopes.Add("thespecialscope")

	// test that a scoped endpoint contains rewrite and auth snippets
	annotations = grp.getAnnotations("service")
	if snippet, ok = annotations["ingress.kubernetes.io/configuration-snippet"]; !ok {
		t.Error("Error, annotations missing configuration-snippet")
	}
	if !strings.Contains(snippet, "auth_request") {
		t.Error("Error, scoped endpoints group doesn't contains auth snippet")
	}
	if !strings.Contains(snippet, "thespecialscope") {
		t.Error("Error, misconfigured snippet for scoped endpoints group")
	}
	if !strings.Contains(snippet, "rewrite") {
		t.Error("Error, expected rewrite snippet")
	}
}

func TestGetIngress(t *testing.T) {

	grp := endpointGroup{}

	// Test for an endpoint group with no endpoints
	ingress := grp.GetIngress("name", 80)
	if _, ok := ingress.Annotations["ingress.kubernetes.io/configuration-snippet"]; !ok {
		t.Error("Error, Ingress annotations missing configuration-snippet")
	}
	if len(ingress.Spec.Rules[0].HTTP.Paths) != 0 {
		t.Error("Error, expected Ingress to have no paths for a pathless endpoint group")
	}

	grp.append("path1")
	grp.append("path2/other")
	grp.append("/path3/other")

	// Test no annotations are missing
	ingress = grp.GetIngress("name", 80)
	if _, ok := ingress.Annotations["ingress.kubernetes.io/configuration-snippet"]; !ok {
		t.Error("Error, Ingress annotations missing configuration-snippet")
	}

	if len(ingress.Spec.Rules[0].HTTP.Paths) != len(grp.paths) {
		t.Errorf("Error, expected Ingress to have %v paths, found: %v", len(grp.paths), len(ingress.Spec.Rules[0].HTTP.Paths))
	}

	path1 := ingress.Spec.Rules[0].HTTP.Paths[0]
	path2 := ingress.Spec.Rules[0].HTTP.Paths[1]
	path3 := ingress.Spec.Rules[0].HTTP.Paths[2]

	// Test the services are properly configured
	assertService := func(backend extensionsv1beta1.IngressBackend) {
		if backend.ServiceName != "name" {
			t.Errorf("Error, unexpected Ingress service name %v", backend.ServiceName)
		}
		if backend.ServicePort.IntVal != 80 {
			t.Errorf("Error, unexpected Ingress port %v", backend.ServicePort.IntVal)
		}
	}
	assertService(path1.Backend)
	assertService(path2.Backend)
	assertService(path3.Backend)

	// Test the paths are properly formed
	if path1.Path != "/name/path1" {
		t.Errorf("Error, unexpected Ingress path %v", path1.Path)
	}
	if path2.Path != "/name/path2/other" {
		t.Errorf("Error, unexpected Ingress path %v", path2.Path)
	}
	if path3.Path != "/name/path3/other" {
		t.Errorf("Error, unexpected Ingress path %v", path3.Path)
	}
}

const TestConfigFile = `
endpoints:
  - path: /p1
    scopes: [123, abc]
  - path: /path2
  - path: /nonprotected1/other
    scopes: [user]
  - path: two/parts
    scopes: [user]
  - path: /one/
    scopes: [user, admin]
  - path: /nonprotected2/
`

// Test config file unmarshalling
func TestUnmarshallConfig(t *testing.T) {

	config, err := unmarshallConfig([]byte(TestConfigFile))
	if err != nil {
		t.Errorf("Error during config file unmarshalling: %v", err)
	}
	if len(config.Endpoints) != 6 {
		t.Errorf("Error, unexpected number of config paths: %v", config)
	}
}

// Test parsing the config file into a set of Ingress configurations
func TestConfigParsing(t *testing.T) {

	ingresses, err := ParseConfig("name", 80, []byte(TestConfigFile))
	if err != nil {
		t.Errorf("Error during config file parsing: %v", err)
	}

	// Test that there are the expected number of scopes
	if len(ingresses) != 4 {
		t.Errorf("Error, unexpected number of ingress configurations: %v", ingresses)
	}

	// Test that there is a group for each set of scopes
	containsGroup := func (scopes ...string) {
		grp := endpointGroup{}
		grp.scopes = set.New()
		for _, scope := range scopes {
			grp.scopes.Add(scope)
		}

		found := false
		for _, ingress := range ingresses {
			if ingress.Labels["identifier"] == strconv.Itoa(grp.getHashCode()) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Error, Ingress configuration not found for scopes: %v", grp.scopes)
		}
	}

	containsGroup("abc", "123")
	containsGroup("")
	containsGroup("user")
	containsGroup("user", "admin")
}