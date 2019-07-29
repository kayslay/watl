package store

import "github.com/globalsign/mgo/bson"

type Message struct {
	ClientID string `json:"client_id" bson:"client_id"`
	Message  string `json:"message" bson:"message"`
	Index    int    `json:"index" bson:"index"`
}

// EditMessage sets the default bot message
func (s MgoStore) EditMessage(m Message) error {
	c, closer := s.db("messages")
	defer closer.Close()
	// TODO aws kms
	_, err := c.Upsert(bson.M{"client_id": m.ClientID, "index": m.Index}, &m)
	return err

}

// GetMessage gets the default bot reply message
func (s MgoStore) GetMessage(clientID string) (Message, error) {
	var m = Message{}
	c, closer := s.db("messages")
	defer closer.Close()

	err := c.Find(bson.M{"client_id": clientID}).One(&m)
	return m, err
}
