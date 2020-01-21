package actions

import (
	"crypto/tls"
	"encoding/xml"
	"log"
	"os"
	//"os"
)
import "amfxmpp/modules"

type XMLPresence struct {
	XMLName  xml.Name `xml:"presence"`
	Type     string   `xml:"type,omitempty"`
	Id       string   `xml:"id,attr,omitempty"`
	From     string   `xml:"from,attr,omitempty"`
	Priority Node     `xml:"priority,omitempty"`
	Show     Node     `xml:"show,omitempty"`
	C        Node     `xml:"c,omitempty"`
	Status   Node     `xml:"status,omitempty"`
}

func ActionPresence(s string, conn *tls.Conn, user *modules.User) bool {
	log.Println("PRESENSE!")

	var inData = []byte(s)
	data := &XMLPresence{}
	xml.Unmarshal(inData, data)

	status := ""
	if data.Type == "unavailable" {
		status = "offline"
	} else if data.Show.Content == "xa" {
		status = "dnd"
	} else if data.Show.Content == "dnd" {
		status = "dnd"
	} else if data.Show.Content == "away" {
		status = "away"
	} else if data.Show.Content == "invisible" {
		status = "invisible"
	} else {
		status = "active"
	}

	message := ""
	if !data.Status.IsEmpty() {
		message = data.Status.Content
	}

	_, err := modules.DB.Exec(`UPDATE xmpp_roster set contact_state=?, contact_state_date=NOW(), contact_status_message=? where user_id=?`,
		status,
		message,
		user.ID,
	)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	return true
}
