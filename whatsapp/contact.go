package whatsapp

import (
	"github/kayslay/watl/store"
	"strings"

	whatzapp "github.com/Rhymen/go-whatsapp"
)

func loadContact(h *Handler) (map[string]store.Contact, error) {
	cc, err := h.store.GetContact(h.id)
	if err != nil {
		return nil, err
	}
	cl := map[string]store.Contact{}

	for _, c := range cc {
		cl[c.Name] = c
	}

	return cl, nil
}

func addContact(h *Handler, msg whatzapp.TextMessage, blacklist bool) error {

	contact := store.Contact{
		ClientID:    h.id,
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

func deleteContact(h *Handler, msg whatzapp.TextMessage) error {
	return h.store.DeleteContact(h.id, msg.Info.RemoteJid)
}
