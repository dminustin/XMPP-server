package actions

import (
	"amfxmpp/config"
	"fmt"
)

func (a *ActionTemplate) Push_DoDisco() bool {

	DoRespond(a.conn,
		fmt.Sprintf("<iq type=\"result\" xmlns=\"jabber:client\" from=\"push.%s\" id=\"%s\" to=\"%s\">"+
			"<query xmlns=\"http://jabber.org/protocol/disco#info\">"+
			"<identity type=\"push\" category=\"pubsub\" name=\"Push Notifications component\" />"+
			"<feature var=\"tigase:messenger:apns:1\" />"+
			"<feature var=\"fcm-xmpp-api\" />"+
			"<feature var=\"http://jabber.org/protocol/commands\" />"+
			"<feature var=\"urn:xmpp:push:0\" />"+
			"<x type=\"result\" xmlns=\"jabber:x:data\">"+
			"<field type=\"hidden\" var=\"FORM_TYPE\">"+
			"<value>http://jabber.org/network/serverinfo</value>"+
			"</field>"+
			a.GetAbuseInfo()+
			"</x>"+
			"</query>"+
			"</iq>"+
			config.Config.Server.Domain,
			a.data.Id,
			a.user.FullAddr,
		), a.data.Id)
	return true
}
