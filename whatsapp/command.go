package whatsapp

import (
	"errors"
	"fmt"
	"github/kayslay/watl/store"
	"log"
	"strings"

	whatzapp "github.com/Rhymen/go-whatsapp"
)

type storer interface {
	AddContact(store.Contact) error
	GetContact(clientID string) ([]store.Contact, error)
	DeleteContact(clientID, name string) error
	EditMessage(store.Message) error
	GetMessage(clientID string) (store.Message, error)
}

// Command handles commands passed to the system
// the command must start the message. they must prefix the message if they
// are to be executed
func (h *Handler) Command(message whatzapp.TextMessage) (string, error) {
	switch {
	case strings.HasPrefix(message.Text, "#!loki"):
		return "you one of the bad guys now", h.addContact(message, true)
	case strings.HasPrefix(message.Text, "#!thor"):
		return "whitelisted", h.addContact(message, false)
	case strings.HasPrefix(message.Text, "#!odin"):
		return toggleState(h)
	case strings.HasPrefix(message.Text, "#!freyja"):
		return helpText(h, message)
	case strings.HasPrefix(message.Text, "#!sif"):
		return editMessage(h, message)
	case strings.HasPrefix(message.Text, "#!hella"):
		return "normal citizen now", h.deleteContact(message)
	}
	return "", errors.New("unkown command")
}

// toggleState toggle the state between RUNNING and IDLE
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

// helpText returns the command help
func helpText(h *Handler, message whatzapp.TextMessage) (string, error) {
	token := strings.SplitAfterN(strings.TrimSpace(message.Text), " ", 2)
	if len(token) == 2 {
		if v, ok := descriptions[token[1]]; ok {
			return fmt.Sprintf("*%s* %s", token[1], v.longDescription), nil
		}
	}
	str := "\n*COMMAND LIST* \n\n"
	for k, v := range descriptions {
		str += fmt.Sprintf("_*%s* %s_ \n\n", k, v.shortDescription)
	}
	str += "To get more info on a command enter *#!freyja [COMMAND]*.\n\nExample *'#!freyja #!freyja'*"
	return strings.TrimSpace(str), nil
}

// editMessage edit the default message the bot returns
func editMessage(h *Handler, message whatzapp.TextMessage) (string, error) {
	txt := strings.TrimPrefix(message.Text, "#!sif ")
	msg := store.Message{
		ClientID: h.c.Info.Wid,
		Message:  txt,
	}
	err := h.store.EditMessage(msg)
	if err != nil {
		return "", err
	}
	h.message = txt
	return "edited bot's reply message", nil
}

func parseMessage(msg string, str ...string) string {
	vars := []string{"{name}"}
	if len(vars) < len(str) {
		return ""
	}

	for i, v := range str {
		msg = strings.Replace(msg, vars[i], v, -1)
	}
	return msg
}
