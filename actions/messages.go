package actions

import (
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"strings"

	"amfxmpp/config"
	"amfxmpp/modules"
	"amfxmpp/utils"
)

type Message struct {
	XMLName xml.Name `xml:"message"`
	Type    string   `xml:"type,attr,omitempty"`
	Id      string   `xml:"id,attr,omitempty"`
	To      string   `xml:"to,attr,omitempty"`
	From    string   `xml:"from,attr,omitempty"`

	Received Node `xml:"received,omitempty"`
	Body     Node `xml:"body,omitempty"`
	OriginID Node `xml:"origin-id,omitempty"`
	Request  Node `xml:"request,omitempty"`
	Thread   Node `xml:"thread,omitempty"`
}

func ActionMessage(s string, conn *tls.Conn, user *modules.User) bool {

	//log.Printf("[MESSAGE] %s", s)

	var inData = []byte(s)
	data := &Message{}
	xml.Unmarshal(inData, data)

	//fmt.Println(data.Received, s)
	if !data.Received.IsEmpty() {
		//todo implement received
		return true
	}

	if len(data.From) > 0 {
		user.ChangeResource(data.From)
	}

	if data.Body.IsEmpty() {
		//todo Do something
		return true
	}

	tmp := strings.Split(data.To, "@")
	to := tmp[0]
	srv := tmp[1]
	//fmt.Println(srv, to)
	if srv != config.Config.Server.Domain {
		return false
	}
	//todo check for empty message
	//todo check if user is banned
	//todo check if user is blacklisted
	content := data.Body.Content

	find := "https://" + config.Config.Server.Domain + ":" + config.Config.FileServer.DownloadPort + config.Config.FileServer.DownloadPath

	isUploaded := strings.Contains(content, find)

	splittedUplUrl := strings.Split(content, "/")
	content = utils.QuoteText(content)

	if isUploaded {
		content = "Uploaded file"
	}
	res, err := modules.DB.Exec("INSERT INTO messages SET date_create=NOW(), from_user=? , to_user=? , message = ?", user.ID, to, content)
	message_id := int64(0)
	aff, err := res.RowsAffected()
	if (err == nil) && (aff > 0) {
		ins, err := res.LastInsertId()
		if err == nil {
			message_id = ins
		} else {
			log.Println(err)
			return false
		}
	}

	user.LastSentMessageID = message_id

	if isUploaded {
		_, err = modules.DB.Exec(`insert into messages_attachements (message_id, filename, from_id, to_id) select 
		? as message_id, xmpp_uploads.filename, ? as from_id, ? as to_id
		from xmpp_uploads
		where hash=? and from_id=?
		`, message_id, user.ID, to, splittedUplUrl[len(splittedUplUrl)-1], user.ID)

		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
	}

	//log.Println(s, data, err, message_id)
	//os.Exit(1)
	return true
}

func (a *ActionTemplate) Messages_DoDisco() bool {
	//todo implement Upload Max File Size into config
	a.user.DoRespond(a.conn,
		fmt.Sprintf(
			"<iq type=\"result\" xmlns=\"jabber:client\" from=\"message-router@%s\" id=\"%s\" to=\"%s\">"+
				"<query xmlns=\"http://jabber.org/protocol/disco#info\">"+
				"<identity type=\"router\" category=\"component\" name=\"AMF XMPP Server\" />"+
				"<identity type=\"im\" category=\"server\" name=\"AMF XMPP Server\" />"+
				"<feature var=\"http://jabber.org/protocol/commands\" />"+
				"<x type=\"result\" xmlns=\"jabber:x:data\">"+
				"<field type=\"hidden\" var=\"FORM_TYPE\">"+
				"<value>http://jabber.org/network/serverinfo</value>"+
				"</field>"+
				a.GetAbuseInfo()+
				"</x>"+
				"</query>"+
				"</iq>",

			config.Config.Server.Domain,
			a.data.Id,
			a.user.FullAddr,
		), a.data.Id)

	return true
}

func (a *ActionTemplate) ActionGetMessagesArchive() bool {

	a.GetResultHeader()
	//os.Exit(1)

	return true
}
