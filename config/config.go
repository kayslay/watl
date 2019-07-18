package config

import (
	"github.com/globalsign/mgo"
)

// Option ...
type Option func(*Config)

// Config ...
type Config struct {
	session *mgo.Session

	closeFns []func()
}

var config *Config

//New create a singleton config struct
func New(fn ...func(*Config)) *Config {
	if config != nil {
		return config
	}
	config = new(Config)
	for _, v := range fn {
		v(config)
	}
	return config
}

// Close close all connections
func Close() {
	for _, fn := range config.closeFns {
		fn()
	}
}
