package store

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

type Session struct {
	ClientID   string    `json:"client_id" bson:"client_id"`
	UniqueCode string    `json:"unique_code" bson:"unique_code"`
	Session    []byte    `json:"session" bson:"session"`
	CreatedAt  time.Time `json:"created_at" bson:"created_at"`
}

// SaveSession saves a single session
func (s MgoStore) SaveSession(sess Session) error {

	c, closer := s.db("sessions")
	defer closer.Close()

	_, err := c.Upsert(bson.M{"client_id": sess.ClientID}, &sess)
	return err
}

// GetSessions get the saved sessions
func (s MgoStore) GetSessions() ([]Session, error) {
	var ss = []Session{}
	c, closer := s.db("sessions")
	defer closer.Close()

	err := c.Find(bson.M{}).All(&ss)
	return ss, err
}

// DeleteSession deletes a  session
func (s MgoStore) DeleteSession(clientID string) ([]Session, error) {
	var ss = []Session{}
	c, closer := s.db("sessions")
	defer closer.Close()

	err := c.Remove(bson.M{"client_id": clientID})
	return ss, err
}
