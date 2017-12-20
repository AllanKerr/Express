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
	"time"
)

type DataStoreAdapter struct {
	ds core.DataStore
}

const NO_TTL = time.Duration(0)

func NewDatastoreAdapter(ds core.DataStore) *DataStoreAdapter {
	adapter := new(DataStoreAdapter)
	adapter.ds = ds
	return adapter
}

// createSession creates a new session record in the Cassandra database. To avoid revealing credentials
// the form data is not persisted. The sig is the unique session identifier and the request consists
// of the data to be persisted. To handle data that expires the ttl or time-to-live can be used. If the
// data does not expire use NO_TTL. ttl must be non-negative.
//
func (adapter *DataStoreAdapter) createSession(sig string, req fosite.Requester, ttl time.Duration) error {

	logrus.WithFields(logrus.Fields{
		"signature":sig,
		"ttl":ttl,
	}).Debug("creating session")

	session, ok := adapter.ds.GetSession().(*gocql.Session)
	if !ok {
		return errors.New("unexpected session type")
	}

	ses, err := NewSession(sig, req)
	if err != nil {
		return fosite.ErrServerError
	}
	if ttl < 0 {
		logrus.Error("attempted to create session with invalid ttl: " + ttl.String())
		return fosite.ErrServerError
	}

	cols := qb.Insert("default.sessions").
		Columns("signature", "request_id", "requested_at", "client_id", "scopes", "granted_scopes", "session_data")
	if ttl != NO_TTL {
		cols = cols.TTL()
	}
	stmt, names := cols.ToCql()

	q := gocqlx.Query(session.Query(stmt), names)
	if ttl == NO_TTL {
		q = q.BindStruct(ses)
	} else {
		q = q.BindStructMap(ses, qb.M{
			"_ttl": qb.TTL(ttl),
		})
	}
	logrus.Debug("query: " + q.String())
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
		if err == gocql.ErrNotFound {
			return nil, fosite.ErrNotFound
		} else {
			return nil, err
		}
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
		if err == gocql.ErrNotFound {
			return nil, fosite.ErrNotFound
		} else {
			return nil, err
		}
	}

	for _, grant := range c.GrantTypes {
		logrus.Info(grant)
	}
	return &c, nil
}

func (adapter *DataStoreAdapter) CreateAuthorizeCodeSession(ctx context.Context, code string, req fosite.Requester) error {

	logrus.Info("CreateAuthorizeCodeSession")
	return adapter.createSession(code, req, NO_TTL)
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

	logrus.Info("CreateAccessTokenSession: ")
	ttl := req.GetSession().GetExpiresAt(fosite.AccessToken).Sub(time.Now().UTC())
	return adapter.createSession(signature, req, ttl)
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
	return adapter.createSession(signature, req, NO_TTL)
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

	logrus.Debug("RevokeRefreshToken")
	return adapter.revokeToken(requestID)
}

func (adapter *DataStoreAdapter) RevokeAccessToken(ctx context.Context, requestID string) error {

	logrus.Debug("RevokeAccessToken")
	return adapter.revokeToken(requestID)
}

func (adapter *DataStoreAdapter) revokeToken(requestID string) error {

	session, ok := adapter.ds.GetSession().(*gocql.Session)
	if !ok {
		return errors.New("unexpected session type")
	}

	// Get the signatures for the tokens with requestID
	stmt, names := qb.Select("default.sessions").
		Columns("signature").
		Where(qb.Eq("request_id")).
		ToCql()

	// Execute select request
	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{
		"request_id": requestID,
	})
	var sigs []string
	if err := gocqlx.Select(&sigs, q.Query); err != nil {
		logrus.Error(err)
		return fosite.ErrServerError
	}
	q.Release()

	// Check that there are tokens to return
	if len(sigs) == 0 {
		return nil
	}

	// Delete the tokens with the matching signatures
	stmt, names = qb.Delete("default.sessions").
		Where(qb.In("signature")).
		ToCql()

	// Execute delete request
	q = gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{
		"signature": sigs,
	})
	if err := q.ExecRelease(); err != nil {
		logrus.WithFields(logrus.Fields{
			"error" :err,
			"signatures" : sigs,
		}).Error("failed to delete tokens")
		return err
	}
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
