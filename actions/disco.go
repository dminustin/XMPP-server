package actions

import (
	"fmt"

	"amfxmpp/config"
)

func (a *ActionTemplate) ActionDiscoInfo() bool {

	if a.data.To == "upload."+config.Config.Server.Domain {
		return a.FileUpload_DoDisco()
	}
	if a.data.To == "message-router."+config.Config.Server.Domain {
		return a.Messages_DoDisco()
	}
	if a.data.To == "pubsub."+config.Config.Server.Domain {
		return a.Pubsub_DoDisco()
	}
	if a.data.To == "push."+config.Config.Server.Domain {
		return a.Push_DoDisco()
	}

	a.user.DoRespond(a.conn,
		fmt.Sprintf("<iq id='%s' type='result' to='%s' from='%s'>"+
			"<query xmlns=\"http://jabber.org/protocol/disco#info\">"+
			"<identity category=\"account\" type=\"registered\" />"+
			"<identity name=\"PubSub acs-clustered\" category=\"pubsub\" type=\"service\" />"+
			"<feature var=\"http://jabber.org/protocol/pubsub#retrieve-default\" />"+
			"<feature var=\"http://jabber.org/protocol/pubsub#purge-nodes\" />"+
			"<feature var=\"http://jabber.org/protocol/pubsub#subscribe\" />"+
			"<feature var=\"http://jabber.org/protocol/pubsub#member-affiliation\" />"+
			"<feature var=\"http://jabber.org/protocol/pubsub#subscription-notifications\" />"+
			"<feature var=\"http://jabber.org/protocol/pubsub#create-nodes\" />"+
			"<feature var=\"http://jabber.org/protocol/pubsub#outcast-affiliation\" />"+
			"<feature var=\"http://jabber.org/protocol/pubsub#get-pending\" />"+
			"<feature var=\"http://jabber.org/protocol/pubsub#presence-notifications\" />"+
			"<feature var=\"urn:xmpp:ping\" />"+
			"<feature var=\"http://jabber.org/protocol/pubsub#delete-nodes\" />"+
			"<feature var=\"http://jabber.org/protocol/pubsub#config-node\" />"+
			"<feature var=\"http://jabber.org/protocol/pubsub#retrieve-items\" />"+
			"<feature var=\"http://jabber.org/protocol/pubsub#access-whitelist\" />"+
			"<feature var=\"http://jabber.org/protocol/pubsub#access-presence\" />"+
			"<feature var=\"urn:xmpp:mam:1\" />"+
			"<feature var=\"http://jabber.org/protocol/pubsub#instant-nodes\" />"+
			"<feature var=\"http://jabber.org/protocol/pubsub#modify-affiliations\" />"+
			"<feature var=\"http://jabber.org/protocol/pubsub#multi-collection\" />"+
			"<feature var=\"http://jabber.org/protocol/pubsub#create-and-configure\" />"+
			"<feature var=\"http://jabber.org/protocol/pubsub#publisher-affiliation\" />"+
			"<feature var=\"http://jabber.org/protocol/pubsub#access-open\" />"+
			"<feature var=\"http://jabber.org/protocol/pubsub#retrieve-affiliations\" />"+
			"<feature var=\"http://jabber.org/protocol/pubsub#access-authorize\" />"+
			"<feature var=\"jabber:iq:version\" />"+
			"<feature var=\"http://jabber.org/protocol/pubsub#retract-items\" />"+
			"<feature var=\"http://jabber.org/protocol/pubsub#manage-subscriptions\" />"+
			"<feature var=\"tigase:pubsub:1\" />"+
			"<feature var=\"http://jabber.org/protocol/commands\" />"+
			"<feature var=\"http://jabber.org/protocol/pubsub#auto-subscribe\" />"+
			"<feature var=\"http://jabber.org/protocol/pubsub#publish-options\" />"+
			"<feature var=\"http://jabber.org/protocol/pubsub#access-roster\" />"+
			"<feature var=\"http://jabber.org/protocol/pubsub#publish\" />"+
			"<feature var=\"http://jabber.org/protocol/pubsub#collections\" />"+
			"<feature var=\"http://jabber.org/protocol/pubsub#retrieve-subscriptions\" />"+
			"<x xmlns=\"jabber:x:data\" type=\"result\">"+
			"<field var=\"FORM_TYPE\" type=\"hidden\">"+
			"<value>http://jabber.org/network/serverinfo</value>"+
			"</field>"+
			a.GetAbuseInfo()+
			"</x>"+
			"<feature var=\"http://jabber.org/protocol/pubsub#auto-create\" />"+
			"<feature var=\"http://jabber.org/protocol/pubsub#auto-subscribe\" />"+
			"<feature var=\"urn:xmpp:carbons:2\" />"+
			"<feature var=\"http://jabber.org/protocol/stats\" />"+
			"<feature var=\"vcard-temp\" />"+
			"<feature var=\"jabber:iq:auth\" />"+
			"<feature var=\"http://jabber.org/protocol/amp\" />"+
			"<feature var=\"msgoffline\" />"+
			"<feature var=\"http://jabber.org/protocol/disco#info\" />"+
			"<feature var=\"http://jabber.org/protocol/disco#items\" />"+
			"<feature var=\"urn:xmpp:blocking\" />"+
			"<feature var=\"urn:xmpp:ping\" />"+
			"<feature var=\"urn:ietf:params:xml:ns:xmpp-sasl\" />"+
			"<feature var=\"http://jabber.org/protocol/pubsub\" />"+
			"<feature var=\"http://jabber.org/protocol/pubsub#owner\" />"+
			"<feature var=\"http://jabber.org/protocol/pubsub#publish\" />"+
			"<identity category=\"pubsub\" type=\"pep\" />"+
			"<feature var=\"urn:xmpp:pep-vcard-conversion:0\" />"+
			"<feature var=\"urn:xmpp:bookmarks-conversion:0\" />"+
			"<feature var=\"urn:xmpp:archive:auto\" />"+
			"<feature var=\"urn:xmpp:archive:manage\" />"+
			"<feature var=\"urn:xmpp:push:0\" />"+
			"<feature var=\"tigase:push:away:0\" />"+
			"<feature var=\"jabber:iq:roster\" />"+
			"<feature var=\"jabber:iq:roster-dynamic\" />"+
			"<feature var=\"urn:xmpp:mam:1\" />"+
			"<feature var=\"jabber:iq:version\" />"+
			"<feature var=\"urn:xmpp:time\" />"+
			"<feature var=\"jabber:iq:privacy\" />"+
			"<feature var=\"urn:ietf:params:xml:ns:xmpp-bind\" />"+
			"<feature var=\"http://jabber.org/protocol/commands\" />"+
			"<feature var=\"urn:ietf:params:xml:ns:vcard-4.0\" />"+
			"<feature var=\"jabber:iq:private\" />"+
			"<feature var=\"urn:ietf:params:xml:ns:xmpp-session\" />"+
			"</query>"+
			"</iq>",

			a.data.Id,
			a.user.FullAddr,
			a.user.Resource,
			config.Config.Server.Domain,
		), a.data.Id)
	return true
}

func (a *ActionTemplate) ActionDiscoItems() bool {

	a.user.DoRespond(a.conn,
		fmt.Sprintf(

			"<iq type=\"result\" xmlns=\"jabber:client\" from=\"%s\" id=\"%s\" to=\"%s\">"+
				"<query xmlns=\"http://jabber.org/protocol/disco#items\">"+
				"<item jid=\"upload.%s\" name=\"HTTP File Upload component\" />"+
				"<item jid=\"push.%s\" name=\"Push Notifications component\" />"+
				"<item jid=\"message-router@%s\" name=\"AMF Xmpp server\" />"+
				"<item jid=\"pubsub.%s\" name=\"PubSub acs-clustered\" />"+
				"</query>"+
				"</iq>",
			config.Config.Server.Domain,
			a.data.Id,
			a.user.FullAddr,
			config.Config.Server.Domain,
			config.Config.Server.Domain,
			config.Config.Server.Domain,
			config.Config.Server.Domain,
		), a.data.Id)
	return true
}
