package whatsapp

import (
	"github/kayslay/watl/store"
	"strings"

	whatzapp "github.com/Rhymen/go-whatsapp"
)

func (h *Handler) LoadContact() error {
	cc, err := h.store.GetContact(h.c.Info.Wid)
	if err != nil {
		return err
	}
	cl := map[string]store.Contact{}

	for _, c := range cc {
		cl[c.Name] = c
	}

	h.contactList = cl

	return nil
}

func (h *Handler) addContact(msg whatzapp.TextMessage, blacklist bool) error {

	// TODO add encryption to saved state
	contact := store.Contact{
		ClientID:    h.c.Info.Wid,
		Phone:       strings.Split(msg.Info.RemoteJid, "@")[0],
		Blacklisted: blacklist,
		Name:        msg.Info.RemoteJid,
		Action:      "",
	}

	err := h.store.AddContact(contact)
	if err != nil {
		return err
	}

	h.contactList[contact.Name] = contact

	return nil
}

func (h *Handler) deleteContact(msg whatzapp.TextMessage) error {
	return h.store.DeleteContact(h.c.Info.Wid, msg.Info.RemoteJid)
}
