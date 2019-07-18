package whatsapp

import (
	"encoding/json"
	"fmt"
	"github/kayslay/watl/store"
	"io"
	"log"
	"strings"
	"time"

	whatzapp "github.com/Rhymen/go-whatsapp"
	"rsc.io/quote"
)

type Handler struct {
	id          string
	c           *whatzapp.Conn
	w           io.WriteCloser
	state       string
	prevState   string
	initTime    time.Time
	contactList map[string]store.Contact
	store       storer
}

func NewHandler(c *whatzapp.Conn, w io.WriteCloser, s storer) (*Handler, error) {
	var err error
	h := &Handler{
		id:        "cli_client",
		c:         c,
		w:         w,
		store:     s,
		state:     "IDLE",
		initTime:  time.Now(),
		prevState: "",
	}
	h.contactList, err = loadContact(h)

	return h, err
}

func (h *Handler) SetClientID(id string) {
	h.id = id
}

//HandleError needs to be implemented to be a valid WhatsApp handler
func (h *Handler) HandleError(err error) {

	if e, ok := err.(*whatzapp.ErrConnectionFailed); ok {
		log.Printf("Connection failed, underlying error: %v", e.Err)
		// try to restart reccursively
		h.restart()
	} else {
		log.Printf("error occurred: %v\n", err)
	}
}

func (h *Handler) restart() {
	if h.state == "RESTARTING" {
		return
	}
	h.prevState, h.state = h.state, "RESTARTING"

	for {
		log.Println("Waiting 30sec...")
		<-time.After(30 * time.Second)
		log.Println("Reconnecting...")
		err := h.c.Restore()
		h.initTime = time.Now()
		// getout of loop when it has successfully connected
		if err == nil {
			h.prevState, h.state = h.state, h.prevState
			return
		}
	}
}

//Optional to be implemented. Implement HandleXXXMessage for the types you need.
func (h *Handler) HandleTextMessage(message whatzapp.TextMessage) {
	if !time.Unix(int64(message.Info.Timestamp), 0).After(h.initTime) {
		return
	}

	// return
	if message.Info.FromMe {
		// check if it is a command
		if strings.HasPrefix(message.Text, "#!") {
			err := h.command(message)
			if err != nil {
				log.Println("error", err)
				msg := whatzapp.TextMessage{
					Info: whatzapp.MessageInfo{
						RemoteJid: message.Info.RemoteJid,
					},
					Text: "_could not run the ducking command lord master kay!_ 🤖",
				}
				h.c.Send(msg)
				msg.Text = fmt.Sprintf("_%s_ 🤖", err.Error())
				h.c.Send(msg)
				return
			}

			msg := whatzapp.TextMessage{
				Info: whatzapp.MessageInfo{
					RemoteJid: message.Info.RemoteJid,
				},
				Text: "_command active *master kay*_ 🤖",
			}

			h.c.Send(msg)
		}
		return
	}

	// check if it white listed
	c, ok := h.contactList[message.Info.RemoteJid]
	if ok && !c.Blacklisted {
		// it is white listed
		// check actions. perform some whitelist actions
		return
	}

	// check if it is running or contact is blacklisted
	if !(h.state == "RUNNING" || c.Blacklisted) {
		// it is not running and user is not blacklisted
		return
	}

	fmt.Println((message.Text), message.Info.RemoteJid)

	switch {
	case "status@broadcast" == message.Info.RemoteJid:
		h.StatusListener(message)
		fmt.Println("reply status")
		// check selected contact list
	default:
		h.echoMessage(message)
	}
}

func (h *Handler) StatusListener(message whatzapp.TextMessage) {

	if !strings.Contains(message.Text, "@kay") || message.Info.Source.Participant == nil {
		return
	}
	h.sendTofile(message)
	msg := whatzapp.TextMessage{
		Info: whatzapp.MessageInfo{
			RemoteJid: *message.Info.Source.Participant,
		},
		Text: "u summoned me 🤖",
	}
	h.c.Send(msg)

}

func (h *Handler) echoMessage(message whatzapp.TextMessage) {
	if !strings.HasPrefix(message.Info.RemoteJid, "234") || !strings.HasSuffix(message.Info.RemoteJid, "s.whatsapp.net") || message.Info.FromMe {
		return
	}
	// h.sendTofile(message)

	msg := whatzapp.TextMessage{
		Info: whatzapp.MessageInfo{
			RemoteJid: message.Info.RemoteJid,
		},
		Text: fmt.Sprintf("_%s_ 🤖", quote.Glass()),
	}

	go func() {
		c, err := h.c.Read(message.Info.RemoteJid, message.Info.Id)
		if err != nil {
			fmt.Println("error reading")
			return
		}
		fmt.Println("read ☑", <-c)
	}()

	h.c.Send(msg)
	msg.Info.QuotedMessageID = ""
	msg.Text = "_*master kay* is busy now. Heleft left whatsapp for enlightment and deeper understanding of life ._ his telegram may be active 🤖"
	h.c.Send(msg)

}

func (h *Handler) processTextMessage(message whatzapp.TextMessage) {
	// name := contact.GetName(strings.Split(message.Info.RemoteJid, "@")[0])

}

func (h *Handler) sendTofile(message interface{}) {
	b, err := json.Marshal(&message)
	if err != nil {
		fmt.Println("file error:", err)
		return
	}

	b = append([]byte(",\n"), b...)
	_, err = h.w.Write(b)
	if err != nil {
		fmt.Println("file error:", err)
		return
	}
}
