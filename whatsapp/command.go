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
}

func (h *Handler) command(message whatzapp.TextMessage) error {
	switch {
	case strings.HasPrefix(message.Text, "#!loki"):
		return addContact(h, message, true)
	case strings.HasPrefix(message.Text, "#!thor"):
		return addContact(h, message, false)
	case strings.HasPrefix(message.Text, "#!odin"):
		return toggleState(h)
	case strings.HasPrefix(message.Text, "#!hella"):
		return errors.New("hella will always fail")
	}
	return nil
}

func toggleState(h *Handler) error {
	if h.state == "RUNNING" {
		h.prevState, h.state = h.state, "IDLE"
	} else if h.state == "IDLE" {
		h.prevState, h.state = h.state, "RUNNING"
	}

	log.Println("ODIN:", h.state, h.prevState)
	return nil
}
