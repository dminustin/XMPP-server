package actions

import (
	appconfig "amfxmpp/config"
	"amfxmpp/modules"
	"crypto/tls"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"log"
	"strings"
	//"log"
)

func TryTuAuth(login string, password string) bool {
	return true
}

type dataFormat_Auth struct {
	Password  string `xml:",chardata" json:"data"`
	Xmlns     string `xml:"xmlns,attr"`
	Mechanism string `xml:"mechanism,attr"`
}

func ActionAuth(s string, conn *tls.Conn, user *modules.User) bool {

	var inData = []byte(s)
	data := &dataFormat_Auth{}
	err := xml.Unmarshal(inData, data)
	if nil != err {
		fmt.Println("Error unmarshalling from XML", err)
	}

	t, err := base64.StdEncoding.DecodeString(data.Password)
	if err != nil {
		DoRespond(conn, "<failure xmlns=\"urn:ietf:params:xml:ns:xmpp-sasl\">"+
			"<not-authorized /><text xml:lang=\"en\">Password not verified</text></failure>")
		return false
	}
	trimed := strings.Trim(string(t), "\x00")
	trimed = strings.Trim(trimed, " ")

	resp := strings.Split(trimed, "\x00")
	log.Println(resp, len(resp))
	if len(resp) != 2 {
		DoRespond(conn, "<failure xmlns=\"urn:ietf:params:xml:ns:xmpp-sasl\">"+
			"<not-authorized /><text xml:lang=\"en\">Password not verified</text></failure>")
		return false
	}
	login := resp[0]
	password := resp[1]
	fmt.Println(login, password)
	if TryTuAuth(login, password) {
		user.Authorized = true
		user.ID = login
		user.UID = login + "@" + appconfig.Config.Server.Domain
		DoRespond(conn, "<success xmlns=\"urn:ietf:params:xml:ns:xmpp-sasl\" />")
		return true
	} else {
		DoRespond(conn, "<failure xmlns=\"urn:ietf:params:xml:ns:xmpp-sasl\">"+
			"<not-authorized /><text xml:lang=\"en\">Password not verified</text></failure>")
		return false
	}
}
