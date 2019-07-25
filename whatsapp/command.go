package whatsapp

import (
	"errors"
	"github/kayslay/watl/store"
	"log"
	"strings"

	whatzapp "github.com/Rhymen/go-whatsapp"
)

type storer interface {
	AddContact(store.Contact) error
	GetContact(clientID string) ([]store.Contact, error)
	DeleteContact(clientID, name string) error
}

func (h *Handler) Command(message whatzapp.TextMessage) (string, error) {
	switch {
	case strings.HasPrefix(message.Text, "#!loki"):
		return "you one of the bad guys now", addContact(h, message, true)
	case strings.HasPrefix(message.Text, "#!thor"):
		return "whitelisted", addContact(h, message, false)
	case strings.HasPrefix(message.Text, "#!odin"):
		return toggleState(h)
	case strings.HasPrefix(message.Text, "#!hella"):
		return "normal citizen now", deleteContact(h, message)
	}
	return "", errors.New("unkown command")
}

func toggleState(h *Handler) (string, error) {
	if h.state == "RUNNING" {
		h.prevState, h.state = h.state, "IDLE"
	} else if h.state == "IDLE" {
		h.prevState, h.state = h.state, "RUNNING"
	}

	log.Println("ODIN:", h.state, h.prevState)
	message := "everybody apart from whitlisted contact has been added to do not disturb"
	if h.state == "IDLE" {
		message = "everybody apart from blacklisted contact can send messages in peace"
	}
	return message, nil
}
