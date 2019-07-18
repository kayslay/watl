package config

import (
	"log"
	"os"
	"time"

	"github.com/globalsign/mgo"
)

type Closer interface {
	Close()
}

// NewMongoConnect create the mongodb connections
func NewMongoConnect(c *Config) {
	var err error
	info := &mgo.DialInfo{
		Addrs:    []string{os.Getenv("MONGO_URI")},
		Timeout:  60 * time.Second,
		Database: os.Getenv("MONGO_URI_DATABASE"),
		Username: os.Getenv("MONGO_USER"),
		Password: os.Getenv("MONGO_PASS"),
	}

	c.session, err = mgo.DialWithInfo(info)
	if err != nil {
		panic(err.Error())
	}

	database := c.session.DB(os.Getenv("MONGO_URI_DATABASE"))
	if err := ensureIndex(database); err != nil {
		panic(err)
	}

	c.closeFns = append(c.closeFns, func() {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
			}
		}()
		log.Println("closing mgo")

		c.session.Close()
	})

}

// Mgo take the database and the collection name to return
func Mgo(c string) (*mgo.Collection, Closer) {
	session := config.session.Copy()
	v := session.DB(os.Getenv("MONGO_URI_DATABASE"))

	return v.C(c), session
}

func ensureIndex(d *mgo.Database) error {
	// TODO add indexes to improve speed
	rIndex := mgo.Index{
		Key:    []string{"client_id", "phone"},
		Unique: true,
	}

	indexes := map[string][]mgo.Index{
		"contacts": {rIndex},
	}

	for i := range indexes {
		for j := range indexes[i] {
			err := d.C(i).EnsureIndex(indexes[i][j])
			if err != nil {
				return err
			}

		}
	}
	return nil
}
