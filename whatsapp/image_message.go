package whatsapp

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	whatzapp "github.com/Rhymen/go-whatsapp"
)

//Example for media handling. Video, Audio, Document are also possible in the same way
func (h *Handler) HandleImageMessage(message whatzapp.ImageMessage) {
	if !time.Unix(int64(message.Info.Timestamp), 0).After(h.initTime) {
		return
	}
	if os.Getenv("ENV") == "production" {
		return
	}
	h.sendTofile(message)
	data, err := message.Download()
	if err != nil {
		return
	}
	filename := fmt.Sprintf("data/img/%v.%v" /* os.TempDir(), */, message.Info.Id, strings.Split(message.Type, "/")[1])
	file, err := os.Create(filename)
	defer file.Close()
	if err != nil {
		return
	}
	_, err = file.Write(data)
	if err != nil {
		return
	}
	log.Printf("%v %v\n\timage reveived, saved at:%v\n", message.Info.Timestamp, message.Info.RemoteJid, filename)
}
