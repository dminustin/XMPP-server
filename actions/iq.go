package actions

import (
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"reflect"

	"amfxmpp/modules"
)

type IQ struct {
	// Info/Query
	XMLName   xml.Name     `xml:"iq"`
	To        string       `xml:"to,attr,omitempty"`
	From      string       `xml:"from,attr,omitempty"`
	Id        string       `xml:"id,attr,omitempty"`
	Type      string       `xml:"type,attr,omitempty"`
	Bind      NestedStruct `xml:"bind,omitempty"`
	Session   NestedStruct `xml:"session,omitempty"`
	Query     NestedStruct `xml:"query,omitempty"`
	Enable    NestedStruct `xml:"enable,omitempty"`
	Request   NestedStruct `xml:"request,omitempty"`
	Blocklist NestedStruct `xml:"blocklist,omitempty"`
	Pubsub    NestedStruct `xml:"pubsub,omitempty"`
	Ping      NestedStruct `xml:"ping,omitempty"`
	Message   NestedStruct `xml:"message,omitempty"`
	VCard     NestedStruct `xml:"vCard,omitempty"`
}

type NestedStruct struct {
	Xmlns   string `xml:"xmlns,attr,omitempty"`
	Payload Node   `xml:"resource,omitempty"`
	Storage Node   `xml:"storage,omitempty"`
	Queryid string `xml:"queryid,omitempty"`
	Roster  Node   `xml:"roster,omitempty"`
	Items   Node   `xml:"items,omitempty"`
	Item    Node   `xml:"item,omitempty"`
	Group   Node   `xml:"group,omitempty"`
	Body    Node   `xml:"body,omitempty"`
	Thread  Node   `xml:"thread,omitempty"`
	Publish Node   `xml:"publish,omitempty"`
}

type Node struct {
	XMLName xml.Name
	Attrs   []xml.Attr `xml:"-,omitempty"`
	Content string     `xml:",cdata"`
	Nodes   []Node     `xml:",any,omitempty"`
	Node    string     `xml:"node,attr,omitempty"`
	Jid     string     `xml:"jid,attr,omitempty"`
	Item    struct {
		Id       string `xml:"id,omitempty"`
		Activity struct {
			Xmlns   string     `xml:"xmlns,attr,omitempty"`
			Attrs   []xml.Attr `xml:"-,omitempty"`
			Content string     `xml:",cdata"`
		} `xml:"activity,omitempty"`
	} `xml:"item,omitempty"`
}

func (s Node) IsEmpty() bool {
	return reflect.DeepEqual(s, Node{})
}
func (s NestedStruct) IsEmpty() bool {
	return reflect.DeepEqual(s, NestedStruct{})
}

func ActionIQ(s string, conn *tls.Conn, user *modules.User) bool {

	var inData = []byte(s)
	data := &IQ{}
	err := xml.Unmarshal(inData, data)
	if nil != err {
		fmt.Println("Error unmarshalling from XML", err, s)
		os.Exit(1)
	}

	if len(data.From) > 0 {
		user.ChangeResource(data.From)
	}

	modules.WriteQueChan(data.Id, s)

	cmd := ""
	log.Println(s)
	//todo clear this dirty code!
	if !data.Bind.Payload.IsEmpty() {
		cmd = "bind"
	} else if !data.Enable.IsEmpty() {
		cmd = "enable"
	} else if !data.Blocklist.IsEmpty() {
		cmd = "blocklist"
	} else if !data.Pubsub.IsEmpty() {
		cmd = "pubsub"
	} else if !data.VCard.IsEmpty() {
		if data.Type == "set" {
			cmd = "vcard.set"
		} else if data.Type == "get" {
			cmd = "vcard.get"
		} else {
			//wtf
			log.Println("Unknown request", s)
		}
	} else if !data.Ping.IsEmpty() {
		cmd = "ping"
	} else if !data.Query.Roster.IsEmpty() {
		cmd = "roster"
	} else if !data.Session.IsEmpty() {
		cmd = "session"
	} else if !data.Query.IsEmpty() {
		if data.Query.Xmlns == "http://jabber.org/protocol/disco#info" {
			cmd = "disco.info"
		} else if data.Query.Xmlns == "http://jabber.org/protocol/disco#items" {
			cmd = "disco.items"
		} else if data.Query.Xmlns == "jabber:iq:roster" && data.Type == "get" {
			cmd = "roster.get"
		} else if data.Query.Xmlns == "jabber:iq:roster" && data.Type == "set" {
			cmd = "roster.set"
		} else if data.Query.Xmlns == "urn:xmpp:mam:1" {
			cmd = "messages.archive"
		} else if !data.Query.Storage.IsEmpty() {
			cmd = "storage"
		}
	} else if !data.Request.IsEmpty() {
		if data.Request.Xmlns == "urn:xmpp:http:upload:0" {
			cmd = "file.upload.request"
		}
	}

	//log.Printf("\n\n%s == %s\n%s\n\n", cmd, data.Query.Xmlns, s)

	var oActionTemplate = ActionTemplate{user: user, conn: conn, data: data}
	switch cmd {
	case "messages.archive":
		{
			return oActionTemplate.ActionGetMessagesArchive()
			break
		}
	case "bind":
		{
			return oActionTemplate.ActionBind()
			break
		}
	case "ping":
		{
			return oActionTemplate.ActionPong()
			break
		}
	case "vcard.set":
		{
			return oActionTemplate.ActionVCardSet()
			break
		}
	case "vcard.get":
		{
			return oActionTemplate.ActionVCardGet()
			break
		}
	case "pubsub":
		{
			return oActionTemplate.ActionPubsub()
			break
		}
	case "blocklist":
		{
			return oActionTemplate.ActionBlockList()
			break
		}
	case "session":
		{
			return oActionTemplate.ActionSession()
			break
		}
	case "storage":
		{
			return oActionTemplate.ActionStorage()
			break
		}
	case "roster":
		{
			return oActionTemplate.ActionRoster()
			break
		}
	case "roster.get":
		{
			return oActionTemplate.ActionRosterGetList()
			break
		}
	case "roster.set":
		{
			return oActionTemplate.ActionRosterSet()
			break
		}
	case "enable":
		{
			return oActionTemplate.ActionEnable()
			break
		}
	case "disco.info":
		{
			return oActionTemplate.ActionDiscoInfo()
			break
		}
	case "disco.items":
		{
			return oActionTemplate.ActionDiscoItems()
			break
		}
	case "file.upload.request":
		{
			return oActionTemplate.ActionRequestFileUpload()
			break
		}
	default:
		{
			type t = struct {
				raw interface{}
				cmd string
				s   string
			}
			(&modules.AppLogStruct{LogType: "ERROR", LogMessage: "Unanswered IQ",
				LogData: t{cmd: cmd, raw: data.Enable, s: s},
			}).WriteAppLog()
			return true
		}
	}

	return true

}
