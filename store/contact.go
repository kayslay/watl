package store

import (
	"github/kayslay/watl/config"

	"github.com/globalsign/mgo"

	"github.com/globalsign/mgo/bson"
)

type Contact struct {
	ID          interface{} `json:"-" bson:"_id"`
	ClientID    string      `json:"client_id" bson:"client_id"`
	Name        string      `json:"name" bson:"name"`
	Phone       string      `json:"phone" bson:"phone"`
	Blacklisted bool        `json:"blacklisted" bson:"blacklisted"`
	Action      string      `json:"action" bson:"action"` //what should happen to it based on where it lies
}

type MgoStore struct {
	db func(string) (*mgo.Collection, config.Closer)
}

func NewMgo() *MgoStore {
	return &MgoStore{db: config.Mgo}
}

func (s MgoStore) AddContact(ct Contact) error {
	// delete contact then add new list
	c, closer := s.db("contacts")
	defer closer.Close()
	// TODO aws kms
	c.Remove(bson.M{"client_id": ct.ClientID, "phone": ct.Phone})
	ct.ID = bson.NewObjectId()
	err := c.Insert(&ct)
	return err
}

func (s MgoStore) GetContact(clientID string) ([]Contact, error) {
	var cc = []Contact{}
	c, closer := s.db("contacts")
	defer closer.Close()

	err := c.Find(bson.M{"client_id": clientID}).All(&cc)
	//
	return cc, err
}

func (s MgoStore) DeleteContact(clientID, name string) error {
	c, closer := s.db("contacts")
	defer closer.Close()
	// TODO aws kms
	return c.Remove(bson.M{"client_id": clientID, "name": name})
}
