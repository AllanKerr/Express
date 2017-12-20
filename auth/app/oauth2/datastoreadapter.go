package oauth2

import (
	"github.com/ory/fosite"
	"context"
	"app/core"
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
	"github.com/gocql/gocql"
	"errors"
	"github.com/sirupsen/logrus"
)

type DataStoreAdapter struct {
	ds core.DataStore
}

func NewDatastoreAdapter(ds core.DataStore) *DataStoreAdapter {
	adapter := new(DataStoreAdapter)
	adapter.ds = ds
	return adapter
}

func (adapter *DataStoreAdapter) createSession(sig string, req fosite.Requester) error {



	session, ok := adapter.ds.GetSession().(*gocql.Session)
	if !ok {
		return errors.New("unexpected session type")
	}
	ses, err := NewSession(sig, req)
	if err != nil {
		return fosite.ErrServerError
	}

	stmt, names := qb.Insert("default.sessions").
		Columns("signature", "request_id", "requested_at", "client_id", "scopes", "granted_scopes", "session_data").
		ToCql()

	q := gocqlx.Query(session.Query(stmt), names).BindStruct(ses)
	if err := q.ExecRelease(); err != nil {
		logrus.WithField("error", err).Error("insert failed")
		return fosite.ErrServerError
	}
	return nil
}

func (adapter *DataStoreAdapter) getSession(sig string) (fosite.Requester, error) {

	session, ok := adapter.ds.GetSession().(*gocql.Session)
	if !ok {
		return nil, errors.New("unexpected session type")
	}
	stmt, names := qb.Select("default.sessions").
		Where(qb.Eq("signature")).
		ToCql()

	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{
		"signature": sig,
	})

	var s Session
	if err := gocqlx.Get(&s, q.Query); err != nil {
		return nil, err
	}
	return &s, nil
}

func (adapter *DataStoreAdapter) CreateClient(client *Client) error {

	session, ok := adapter.ds.GetSession().(*gocql.Session)
	if !ok {
		return errors.New("unexpected session type")
	}
	stmt, names := qb.Insert("default.clients").
		Columns("id", "secret_hash", "redirect_uris", "grant_types", "response_types", "scopes", "public").
		ToCql()

	// bind the new client to be inserted
	q := gocqlx.Query(session.Query(stmt), names).BindStruct(client)
	if err := q.ExecRelease(); err != nil {
		logrus.WithField("error", err).Error("failed to insert client")
		return err
	}
	return nil
}

func (adapter *DataStoreAdapter) GetClient(_ context.Context, id string) (fosite.Client, error) {

	logrus.Info("GetClient")

	session, ok := adapter.ds.GetSession().(*gocql.Session)
	if !ok {
		return nil, errors.New("unexpected session type")
	}
	stmt, names := qb.Select("default.clients").
		Where(qb.Eq("id")).
		ToCql()

	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{
		"id": id,
	})

	var c Client
	if err := gocqlx.Get(&c, q.Query); err != nil {
		return nil, err
	}

	for _, grant := range c.GrantTypes {
		logrus.Info(grant)
	}
	return &c, nil
}

func (adapter *DataStoreAdapter) CreateAuthorizeCodeSession(ctx context.Context, code string, req fosite.Requester) error {

	logrus.Info("CreateAuthorizeCodeSession")
	return adapter.createSession(code, req)
}

func (adapter *DataStoreAdapter) GetAuthorizeCodeSession(ctx context.Context, code string, _ fosite.Session) (fosite.Requester, error) {

	logrus.Info("GetAuthorizeCodeSession")
	return adapter.getSession(code)
}

func (adapter *DataStoreAdapter) DeleteAuthorizeCodeSession(ctx context.Context, code string) error {

	logrus.Info("DeleteAuthorizeCodeSession")

	return nil
}

func (adapter *DataStoreAdapter) CreateAccessTokenSession(ctx context.Context, signature string, req fosite.Requester) error {

	logrus.Info("CreateAccessTokenSession")
	return adapter.createSession(signature, req)
}

func (adapter *DataStoreAdapter) GetAccessTokenSession(ctx context.Context, signature string, _ fosite.Session) (fosite.Requester, error) {

	logrus.Info("GetAccessTokenSession")
	return adapter.getSession(signature)
}

func (adapter *DataStoreAdapter) DeleteAccessTokenSession(ctx context.Context, signature string) error {

	logrus.Info("DeleteAccessTokenSession")

	return nil
}

func (adapter *DataStoreAdapter) CreateRefreshTokenSession(ctx context.Context, signature string, req fosite.Requester) error {

	logrus.Info("CreateRefreshTokenSession")
	return adapter.createSession(signature, req)
}

func (adapter *DataStoreAdapter) GetRefreshTokenSession(ctx context.Context, signature string, _ fosite.Session) (fosite.Requester, error) {

	logrus.Info("GetRefreshTokenSession")
	return adapter.getSession(signature)
}

func (adapter *DataStoreAdapter) DeleteRefreshTokenSession(ctx context.Context, signature string) error {

	logrus.Info("DeleteRefreshTokenSession")

	return nil
}

func (adapter *DataStoreAdapter) Authenticate(ctx context.Context, name string, secret string) error {

	logrus.Info("Authenticate")

	session, ok := adapter.ds.GetSession().(*gocql.Session)
	if !ok {
		return errors.New("unexpected session type")
	}
	stmt, names := qb.Select("default.users").
		Where(qb.Eq("username")).
		ToCql()

	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{
		"username": name,
	})

	var u User
	if err := gocqlx.Get(&u, q.Query); err != nil {
		logrus.Error(err)
		return fosite.ErrRequestUnauthorized
	}
	err := u.VerifyPassword(secret)
	if err != nil {
		logrus.Error(err)
		return fosite.ErrRequestUnauthorized
	}

	return nil
}

func (adapter *DataStoreAdapter) RevokeRefreshToken(ctx context.Context, requestID string) error {

	logrus.Info("RevokeRefreshToken")

	return nil
}

func (adapter *DataStoreAdapter) RevokeAccessToken(ctx context.Context, requestID string) error {

	logrus.Info("RevokeAccessToken")

	return nil
}

func (adapter *DataStoreAdapter) CreateUser(user *User) error {

	session, ok := adapter.ds.GetSession().(*gocql.Session)
	if !ok {
		return errors.New("unexpected session type")
	}

	stmt, names := qb.Insert("default.users").
		Columns("username", "password_hash", "role").
		ToCql()

	q := gocqlx.Query(session.Query(stmt), names).BindStruct(user)
	if err := q.ExecRelease(); err != nil {
		logrus.WithField("error", err).Fatal("failed to create user")
		return err
	}
	return nil
}
