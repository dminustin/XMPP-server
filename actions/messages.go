package actions

import (
	"amfxmpp/config"
	"fmt"
	"os"
)

func (a *ActionTemplate) Messages_DoDisco() bool {
	//todo implement Upload Max File Size into config
	DoRespond(a.conn,
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
	os.Exit(1)

	return true
}
