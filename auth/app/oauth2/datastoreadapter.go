package oauth2

import (
	"github.com/ory/fosite"
	"context"
	"app/core"
)

type DataStoreAdapter struct {
	datastore core.DataStore
}

func NewDatastoreAdapter(ds core.DataStore) *DataStoreAdapter {
	adapter := new(DataStoreAdapter)
	adapter.datastore = ds
	return adapter
}

func (adapter *DataStoreAdapter) GetClient(_ context.Context, id string) (fosite.Client, error) {

	result := &Client{}
	return result, nil
}

func (adapter *DataStoreAdapter) CreateAuthorizeCodeSession(ctx context.Context, code string, req fosite.Requester) error {

	return nil
}

func (adapter *DataStoreAdapter) GetAuthorizeCodeSession(ctx context.Context, code string, _ fosite.Session) (fosite.Requester, error) {

	return nil, nil
}

func (adapter *DataStoreAdapter) DeleteAuthorizeCodeSession(ctx context.Context, code string) error {

	return nil
}

func (adapter *DataStoreAdapter) CreateAccessTokenSession(ctx context.Context, signature string, req fosite.Requester) error {

	return nil
}

func (adapter *DataStoreAdapter) GetAccessTokenSession(ctx context.Context, signature string, _ fosite.Session) (fosite.Requester, error) {

	return nil, nil
}

func (adapter *DataStoreAdapter) DeleteAccessTokenSession(ctx context.Context, signature string) error {

	return nil
}

func (adapter *DataStoreAdapter) CreateRefreshTokenSession(ctx context.Context, signature string, req fosite.Requester) error {

	return nil
}

func (adapter *DataStoreAdapter) GetRefreshTokenSession(ctx context.Context, signature string, _ fosite.Session) (fosite.Requester, error) {

	return nil, nil
}

func (adapter *DataStoreAdapter) DeleteRefreshTokenSession(ctx context.Context, signature string) error {

	return nil
}

func (adapter *DataStoreAdapter) Authenticate(ctx context.Context, name string, secret string) error {

	return nil
}

func (ds *DataStoreAdapter) RevokeRefreshToken(ctx context.Context, requestID string) error {

	return nil
}

func (ds *DataStoreAdapter) RevokeAccessToken(ctx context.Context, requestID string) error {

	return nil
}


func (adapter *DataStoreAdapter) CreateOpenIDConnectSession(ctx context.Context, authorizeCode string, requester fosite.Requester) error {

	return nil
}

func (adapter *DataStoreAdapter) GetOpenIDConnectSession(ctx context.Context, authorizeCode string, requester fosite.Requester) (fosite.Requester, error) {

	return nil, nil
}

func (adapter *DataStoreAdapter) DeleteOpenIDConnectSession(ctx context.Context, authorizeCode string) error {

	return nil
}
