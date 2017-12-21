package oauth2

import (
	"github.com/ory/fosite"
	"context"
	"app/core"
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
	"github.com/gocql/gocql"
	"github.com/sirupsen/logrus"
	"time"
	"github.com/pkg/errors"
)

// used to indicate sessions with no expiry time-to-live
const NoTtl = time.Duration(0)

// DataStoreAdapter handles mapping between the persistence operations
// required by fosite and the Cassandra data store
type DataStoreAdapter struct {
	ds core.DataStore
	hasher fosite.Hasher
}

// Create a new data store adapter to map between
// the fosite persistence operations and the data store ds.
// ds must be a Cassandra data store.
func NewDataStoreAdapter(ds core.DataStore, hasher fosite.Hasher) *DataStoreAdapter {
	return &DataStoreAdapter{
		ds,
		hasher,
	}
}

// getCqlSession gets the data store's session and interprets it
// as a CQL session. This is forced and is fatal if the session
// is not a CQL session.
//
func (adapter *DataStoreAdapter) getCqlSession() (*gocql.Session) {
	session, ok := adapter.ds.GetSession().(*gocql.Session)
	if !ok {
		logrus.Fatal("unexpected session type when getting CQL session")
	}
	return session
}

// createSession creates a new session record in the Cassandra database. To avoid revealing credentials
// the form data is not persisted. The sig is the unique session identifier and the request consists
// of the data to be persisted. To handle data that expires the ttl or time-to-live can be used. If the
// data does not expire use NoTtl. ttl must be non-negative.
//
func (adapter *DataStoreAdapter) createSession(sig string, req fosite.Requester, ttl time.Duration) error {

	logrus.WithFields(logrus.Fields{
		"signature":sig,
		"ttl":ttl,
	}).Debug("creating session")

	// create a new session struct for the requester
	ses, err := NewSession(sig, req)
	if err != nil {
		logrus.WithField("error", err).Error("failed to create new session")
		return errors.WithStack(fosite.ErrServerError)
	}
	// check the time-to-live is non-negative
	if ttl < 0 {
		logrus.WithField("ttl",ttl).Error("attempted to create session with invalid ttl")
		return errors.WithStack(fosite.ErrServerError)
	}

	session := adapter.getCqlSession()

	// build the insert request to insert the session data
	cols := qb.Insert("default.sessions").
		Columns("signature", "request_id", "requested_at", "client_id", "scopes", "granted_scopes", "session_data")
	if ttl != NoTtl {
		cols = cols.TTL()
	}
	stmt, names := cols.ToCql()
	q := gocqlx.Query(session.Query(stmt), names)
	if ttl == NoTtl {
		q = q.BindStruct(ses)
	} else {
		// bind the ttl if it was provided
		q = q.BindStructMap(ses, qb.M{
			"_ttl": qb.TTL(ttl),
		})
	}

	// insert the session data into the data store
	if err := q.ExecRelease(); err != nil {
		logrus.WithField("error", err).Error("failed to insert session")
		return errors.WithStack(fosite.ErrServerError)
	}
	return nil
}

// gets a session with the specified signature from the Cassandra data store
// If no session is found then ErrNotFound is returned
func (adapter *DataStoreAdapter) getSession(sig string) (fosite.Requester, error) {

	logrus.WithField("signature",sig).Debug("getting session")

	session := adapter.getCqlSession()

	// build the select query to get the session
	stmt, names := qb.Select("default.sessions").
		Where(qb.Eq("signature")).
		ToCql()
	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{
		"signature": sig,
	})
	defer q.Release()

	// get the specified session
	var s Session
	if err := gocqlx.Get(&s, q.Query); err != nil {
		if err == gocql.ErrNotFound {
			err = fosite.ErrNotFound
		} else {
			logrus.WithField("error", err).Error("failed to get session")
		}
		return nil, errors.WithStack(err)
	}
	return &s, nil
}

// delete the session with the specified signature from the Cassandra data store
func (adapter *DataStoreAdapter) deleteSession(sig string) error {

	logrus.WithField("signature",sig).Debug("deleting session")

	session := adapter.getCqlSession()

	// build the delete query for the tokens with the matching signature
	stmt, names := qb.Delete("default.sessions").
		Where(qb.Eq("signature")).
		ToCql()
	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{
		"signature": sig,
	})

	// delete the matching session
	if err := q.ExecRelease(); err != nil {
		logrus.WithField("error", err).Error("failed to delete session")
		return errors.WithStack(err)
	}
	return nil
}

// create a new client in the data store
// client must be non-nil
func (adapter *DataStoreAdapter) CreateClient(client *Client) error {

	logrus.WithField("client", client.Id).Debug("create client")
	if client == nil {
		return fosite.ErrInvalidRequest.WithDebug("attempted to create a nil client")
	}

	session := adapter.getCqlSession()

	// Compute the secret hash if it doesn't already exist
	if client.SecretHash == nil {
		hash, err  := adapter.hasher.Hash([]byte(client.Secret))
		if err != nil {
			logrus.WithField("error", err).Error("failed to hash client secret")
			return errors.WithStack(err)
		}
		client.SecretHash = hash
	}

	// build the insert client query
	stmt, names := qb.Insert("default.clients").
		Columns("id", "secret_hash", "redirect_uris", "grant_types", "response_types", "scopes", "public").
		Unique().
		ToCql()
	q := gocqlx.Query(session.Query(stmt), names).BindStruct(client)

	// insert the new client
	if err := q.ExecRelease(); err != nil {
		logrus.WithField("error", err).Error("failed to insert client")
		return  errors.WithStack(err)
	}
	return nil
}

// get an existing client from the data store using its client id
func (adapter *DataStoreAdapter) GetClient(_ context.Context, id string) (fosite.Client, error) {

	logrus.WithField("client", id).Debug("get client")

	session := adapter.getCqlSession()

	// build the select query to get the client with its id
	stmt, names := qb.Select("default.clients").
		Where(qb.Eq("id")).
		ToCql()
	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{
		"id": id,
	})
	defer q.Release()

	// get the client from the data store
	var c Client
	if err := gocqlx.Get(&c, q.Query); err != nil {
		if err == gocql.ErrNotFound {
			err = fosite.ErrNotFound
		} else {
			logrus.WithField("error", err).Error("failed to get client")
		}
		return nil, errors.WithStack(err)
	}
	return &c, nil
}

// update a specific column for an existing client
func (adapter *DataStoreAdapter) UpdateClient(client *Client, column string) error {

	session := adapter.getCqlSession()

	stmt, names := qb.Update("default.clients").
		Set(column).
		Where(qb.Eq("id")).
		ToCql()

	q := gocqlx.Query(session.Query(stmt), names).BindStruct(client)
	if err := q.ExecRelease(); err != nil {
		logrus.WithField("error", err).Error("failed to update client")
		return errors.WithStack(err)
	}
	return nil
}

// create a new access token session in the data store
func (adapter *DataStoreAdapter) CreateAccessTokenSession(ctx context.Context, signature string, req fosite.Requester) error {

	logrus.Debug("create access token session")
	ttl := req.GetSession().GetExpiresAt(fosite.AccessToken).Sub(time.Now().UTC())
	return adapter.createSession(signature, req, ttl)
}

// get an access token session from the data store with the matching signature
func (adapter *DataStoreAdapter) GetAccessTokenSession(ctx context.Context, signature string, _ fosite.Session) (fosite.Requester, error) {

	logrus.Debug("get access token session")
	return adapter.getSession(signature)
}

// delete an access token session from the data store with the matching signature
func (adapter *DataStoreAdapter) DeleteAccessTokenSession(ctx context.Context, signature string) error {

	logrus.Debug("delete access token session")
	return adapter.deleteSession(signature)
}

// create a new refresh token session in the data store
func (adapter *DataStoreAdapter) CreateRefreshTokenSession(ctx context.Context, signature string, req fosite.Requester) error {

	logrus.Debug("create refresh token session")
	return adapter.createSession(signature, req, NoTtl)
}

// get a refresh token session from the data store with the matching signature
func (adapter *DataStoreAdapter) GetRefreshTokenSession(ctx context.Context, signature string, _ fosite.Session) (fosite.Requester, error) {

	logrus.Debug("get refresh token session")
	return adapter.getSession(signature)
}

// delete a refresh token session from the data store with the matching signature
func (adapter *DataStoreAdapter) DeleteRefreshTokenSession(ctx context.Context, signature string) error {

	logrus.Debug("delete refresh token session")
	return adapter.deleteSession(signature)
}

// authenticate a user with the provided username and secret password
// if the authentication succeeds nil is returned
func (adapter *DataStoreAdapter) GetUser(ctx context.Context, name string) (User, error) {

	logrus.WithField("name", name).Debug("authenticate")

	session := adapter.getCqlSession()

	// build the select request to get the user data fromt he data store
	stmt, names := qb.Select("default.users").
		Where(qb.Eq("username")).
		ToCql()
	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{
		"username": name,
	})
	defer q.Release()

	var u DefaultUser
	if err := gocqlx.Get(&u, q.Query); err != nil {
		logrus.WithField("error", err).Error("failed to get user")
		return nil, errors.WithStack(fosite.ErrRequestUnauthorized)
	}
	return &u, nil
}

// revoke all tokens with the specified request id
func (adapter *DataStoreAdapter) RevokeRefreshToken(ctx context.Context, requestID string) error {

	logrus.Debug("revoke refresh token")
	return adapter.revokeToken(requestID)
}

// revoke all tokens with the specified request id
func (adapter *DataStoreAdapter) RevokeAccessToken(ctx context.Context, requestID string) error {

	logrus.Debug("revoke access token")
	return adapter.revokeToken(requestID)
}

// revoke all tokens with the specified request id
func (adapter *DataStoreAdapter) revokeToken(requestID string) error {

	session := adapter.getCqlSession()

	// get the signatures for the tokens with requestID
	stmt, names := qb.Select("default.sessions").
		Columns("signature").
		Where(qb.Eq("request_id")).
		ToCql()
	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{
		"request_id": requestID,
	})

	// execute select request
	var sigs []string
	if err := gocqlx.Select(&sigs, q.Query); err != nil {

		logrus.WithFields(logrus.Fields{
			"error":err,
			"request":requestID,
		}).Error("failed to lookup request id")

		return errors.WithStack(fosite.ErrServerError)
	}
	q.Release()

	// check that there are tokens to return
	if len(sigs) == 0 {
		return nil
	}

	// delete the tokens with the matching signatures
	stmt, names = qb.Delete("default.sessions").
		Where(qb.In("signature")).
		ToCql()

	// execute delete request
	q = gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{
		"signature": sigs,
	})
	if err := q.ExecRelease(); err != nil {
		logrus.WithFields(logrus.Fields{
			"error" :err,
			"signatures" : sigs,
		}).Error("failed to delete tokens")
		return errors.WithStack(err)
	}
	return nil
}

// create a new user in the data store
// user must be non-nil
func (adapter *DataStoreAdapter) CreateUser(user *DefaultUser) error {

	logrus.WithField("client", user.Username).Debug("create client")
	if user == nil {
		return fosite.ErrInvalidRequest.WithDebug("attempted to create a nil user")
	}

	session := adapter.getCqlSession()

	// Compute the secret hash if it doesn't already exist
	if user.PasswordHash == nil {
		hash, err  := adapter.hasher.Hash([]byte(user.Password))
		if err != nil {
			logrus.WithField("error", err).Error("failed to hash user password")
			return errors.WithStack(err)
		}
		user.PasswordHash = hash
	}

	// build the user insert query
	stmt, names := qb.Insert("default.users").
		Columns("username", "password_hash", "scopes").
		Unique().
		ToCql()
	q := gocqlx.Query(session.Query(stmt), names).BindStruct(user)

	// insert the user
	if err := q.ExecRelease(); err != nil {
		logrus.WithField("error", err).Fatal("failed to create user")
		return errors.WithStack(err)
	}
	return nil
}

// update a specific column for an existing user
func (adapter *DataStoreAdapter) UpdateUser(user User, column string) error {

	session := adapter.getCqlSession()

	stmt, names := qb.Update("default.users").
		Set(column).
		Where(qb.Eq("username")).
		ToCql()

	q := gocqlx.Query(session.Query(stmt), names).BindStruct(user)
	if err := q.ExecRelease(); err != nil {
		logrus.WithField("error", err).Error("failed to update user")
		return errors.WithStack(err)
	}
	return nil
}

// unsupported authorization grant operation
func (adapter *DataStoreAdapter) CreateAuthorizeCodeSession(ctx context.Context, code string, req fosite.Requester) error {

	logrus.Error("unsupported: create authorize code session")
	return adapter.createSession(code, req, NoTtl)
}

// unsupported authorization grant operation
func (adapter *DataStoreAdapter) GetAuthorizeCodeSession(ctx context.Context, code string, _ fosite.Session) (fosite.Requester, error) {

	logrus.Error("unsupported: get authorize code session")
	return adapter.getSession(code)
}

// unsupported authorization grant operation
func (adapter *DataStoreAdapter) DeleteAuthorizeCodeSession(ctx context.Context, code string) error {

	logrus.Error("unsupported: delete authorize code session")
	return adapter.deleteSession(code)
}
