package oauth2

import "testing"

func TestDefaultUser_AppendScope(t *testing.T) {

	user := NewUser("username", "password")
	if len(user.GetScopes()) != 2 {
		t.Errorf("Error, unexpected number of scopes: %v", len(user.Scopes))
	}

	// test appending a duplicate scope
	user.AppendScope(user.GetScopes()[0])
	if len(user.GetScopes()) != 2 {
		t.Errorf("Error, duplicate scope added: %v", len(user.Scopes))
	}

	// test appending a new scope
	user.AppendScope("unique")
	if len(user.GetScopes()) != 3 {
		t.Errorf("Error adding scope: %v", "unique")
	}
}
