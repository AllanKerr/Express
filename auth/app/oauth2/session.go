package oauth2

import (
	"time"
	"github.com/ory/fosite"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/url"
	"github.com/pborman/uuid"
)

type Session struct {
	Signature string
	RequestId string
	RequestedAt time.Time
	ClientId string
	Scopes []string
	GrantedScopes []string
	FormData url.Values
	SessionData []byte
}

func NewSession(sig string, req fosite.Requester) (*Session, error) {

	sessionData, err := json.Marshal(req.GetSession())
	if err != nil {
		return nil, err
	}
	return &Session{
		RequestId: req.GetID(),
		Signature: sig,
		RequestedAt: req.GetRequestedAt(),
		ClientId: req.GetClient().GetID(),
		Scopes: []string(req.GetRequestedScopes()),
		GrantedScopes: []string(req.GetGrantedScopes()),
		FormData: req.GetRequestForm(),
		SessionData: sessionData,
	}, nil
}

func (s *Session) GetID() string {
	return s.RequestId
}

func (s *Session) SetID(id string) {
	if s.RequestId == "" {
		logrus.Warning("session id set with empty id")
		s.RequestId = uuid.New()
	} else {
		s.RequestId = id
	}
}

func (s *Session) GetRequestedAt() time.Time {
	return s.RequestedAt
}

func (s *Session) GetClient() fosite.Client {
	return nil
}

func (s *Session) GetRequestedScopes() fosite.Arguments {
	return fosite.Arguments(s.Scopes)
}

func (s *Session) SetRequestedScopes(scopes fosite.Arguments) {

	s.Scopes = nil
	for _, scope := range scopes {
		s.AppendRequestedScope(scope)
	}
}

func (s *Session) AppendRequestedScope(scope string) {

	for _, cur := range s.Scopes {
		if scope == cur {
			return
		}
	}
	s.Scopes = append(s.Scopes, scope)
}

func (s *Session) GetGrantedScopes() fosite.Arguments {
	return fosite.Arguments(s.GrantedScopes)
}

func (s *Session) GrantScope(scope string) {

	for _, cur := range s.GrantedScopes {
		if scope == cur {
			return
		}
	}
	s.GrantedScopes = append(s.GrantedScopes, scope)
}

func (s *Session) GetSession() fosite.Session {

	if s.SessionData == nil {
		return nil
	}
	var session fosite.DefaultSession
	err := json.Unmarshal(s.SessionData, &session)
	if err != nil {
		logrus.WithField("err", err).Error("session get session data un-marshall failed")
	}
	return &session
}

func (s *Session) SetSession(session fosite.Session) {

	var sessionData []byte
	if session == nil {
		sessionData = nil
	} else {
		var err error
		sessionData, err = json.Marshal(session)
		if err != nil {
			logrus.WithField("error", err).Error("session set session data marshall failed")
			return
		}
	}
	s.SessionData = sessionData
}

func (s *Session) GetRequestForm() url.Values {
	return s.FormData
}

func (s *Session)Merge(requester fosite.Requester) {

	for _, scope := range requester.GetRequestedScopes() {
		s.AppendRequestedScope(scope)
	}
	for _, scope := range requester.GetGrantedScopes() {
		s.GrantScope(scope)
	}
	s.RequestedAt = requester.GetRequestedAt()
	s.ClientId = requester.GetClient().GetID()
	s.SetSession(requester.GetSession())

	for k, v := range requester.GetRequestForm() {
		s.FormData[k] = v
	}
}
