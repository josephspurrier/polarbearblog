package websession

import (
	"time"
)

// Encrypter -
type Encrypter interface {
	Encrypt(data []byte) (encrypted []byte, err error)
	Decrypt(data []byte) (decrypted []byte, err error)
}

// JSONSession -
type JSONSession struct {
	sessionstorer Sessionstorer
	encrypter     Encrypter
}

// NewJSONSession -
func NewJSONSession(sd Sessionstorer, encrypter Encrypter) (*JSONSession, error) {
	s := &JSONSession{
		sessionstorer: sd,
		encrypter:     encrypter,
	}

	return s, nil
}

// Find returns the data for a given session token from the store. If the
// session token is not found or is expired, the returned exists flag will be
// set to false.
func (s *JSONSession) Find(token string) (b []byte, exists bool, err error) {
	sd := new(SessionDatabase)
	err = sd.Load(s.sessionstorer, s.encrypter)
	if err != nil {
		return nil, false, err
	}

	record, found := sd.Records[token]
	if !found {
		return nil, false, nil
	}

	return record.Data, true, nil
}

// Commit adds a session token and data to the store with the given expiry time.
// If the session token already exists then the data and expiry time are updated.
func (s *JSONSession) Commit(token string, b []byte, expiry time.Time) error {
	sd := new(SessionDatabase)
	err := sd.Load(s.sessionstorer, s.encrypter)
	if err != nil {
		return err
	}

	sd.Records[token] = SessionData{ID: token, Data: b, Expire: expiry}

	return sd.Save(s.sessionstorer, s.encrypter)
}

// Delete removes a session token and corresponding data from the store.
func (s *JSONSession) Delete(token string) error {
	sd := new(SessionDatabase)
	err := sd.Load(s.sessionstorer, s.encrypter)
	if err != nil {
		return err
	}

	delete(sd.Records, token)

	return sd.Save(s.sessionstorer, s.encrypter)
}
