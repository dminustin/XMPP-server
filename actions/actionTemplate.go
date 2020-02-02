package actions

import (
	"crypto/tls"
	"fmt"
	"time"

	"amfxmpp/config"
	"amfxmpp/modules"
)

type ActionTemplate struct {
	data *IQ
	conn *tls.Conn
	user *modules.User
}

func (*ActionTemplate) GetDateVersion() string {
	dt := time.Now().Unix()
	return fmt.Sprintf("%d", dt)
}

func (*ActionTemplate) GetAbuseInfo() string {
	//todo implement admins
	return fmt.Sprintf("<field var=\"abuse-addresses\" type=\"text-multi\">"+
		"<value>xmpp:1@%s</value>"+
		"</field>", config.Config.Server.Domain)
}

func (a *ActionTemplate) GetResultHeader() string {
	return fmt.Sprintf("<iq type=\"result\" xmlns=\"jabber:client\" id=\"%s\" to=\"%s\">",
		a.data.Id,
		a.user.FullAddr,
	)

}
