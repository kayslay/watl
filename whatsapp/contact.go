package whatsapp

import (
	"github/kayslay/watl/store"
	"strings"

	whatzapp "github.com/Rhymen/go-whatsapp"
)

// LoadContact load the contacts on start
func (h *Handler) loadContact() error {

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

// addContact adds a contact to a whitelist/blacklist
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

// deleteContact deletes a contact from a blacklist/whitelist
func (h *Handler) deleteContact(msg whatzapp.TextMessage) (err error) {
	delete(h.contactList, msg.Info.RemoteJid)
	return h.store.DeleteContact(h.c.Info.Wid, msg.Info.RemoteJid)
}
