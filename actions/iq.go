package actions

import (
	"amfxmpp/modules"
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"log"
	"reflect"
)

type IQ struct { // Info/Query
	XMLName xml.Name     `xml:"iq"`
	To      string       `xml:"to,attr,omitempty"`
	Id      string       `xml:"id,attr"`
	Type    string       `xml:"type,attr"`
	Bind    NestedStruct `xml:"bind"`
	Session NestedStruct `xml:"session"`
	Query   NestedStruct `xml:"query"`
	Enable  NestedStruct `xml:"enable"`
	Request NestedStruct `xml:"request"`
}
type NestedStruct struct {
	Xmlns   string `xml:"xmlns,attr,omitempty"`
	Payload Node   `xml:"resource,omitempty"`
}

type Node struct {
	XMLName xml.Name
	Attrs   []xml.Attr `xml:"-"`
	Content string     `xml:",cdata"`
	Nodes   []Node     `xml:",any"`
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
		fmt.Println("Error unmarshalling from XML", err)
	}
	log.Printf("[IQ %s] %s", data.Type, data)

	cmd := ""

	if !data.Bind.Payload.IsEmpty() {
		cmd = "bind"
	} else if !data.Session.IsEmpty() {
		cmd = "session"
	} else if !data.Query.IsEmpty() {
		if data.Query.Xmlns == "http://jabber.org/protocol/disco#info" {
			cmd = "disco.info"
		} else if data.Query.Xmlns == "http://jabber.org/protocol/disco#items" {
			cmd = "disco.items"
		}
	} else if !data.Request.IsEmpty() {
		if data.Request.Xmlns == "urn:xmpp:http:upload:0" {
			cmd = "file.upload.request"
		}
	}

	//log.Printf("\n\n%s == %s\n\n", cmd, data.Request.Xmlns)

	var oActionTemplate = ActionTemplate{user: user, conn: conn, data: data}
	log.Printf("[XXXIQ %s] %s", cmd, data)
	switch cmd {
	case "bind":
		{
			return oActionTemplate.ActionBind()
			break
		}
	case "session":
		{
			return oActionTemplate.ActionSession()
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
			return true
		}
	}

	log.Printf("[IQ %s] %s", data.Type, cmd)
	return true

}
