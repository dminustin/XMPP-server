package actions

import (
	"amfxmpp/config"
	"amfxmpp/modules"
	"crypto/tls"
	"fmt"
)

type ActionTemplate struct {
	data *IQ
	conn *tls.Conn
	user *modules.User
}

func (*ActionTemplate) GetAbuseInfo() string {
	//todo implement admins
	return fmt.Sprintf("<field var=\"abuse-addresses\" type=\"text-multi\">"+
		"<value>xmpp:1@%s</value>"+
		"</field>", config.Config.Server.Domain)
}
