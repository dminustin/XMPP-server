package actions

import (
	"amfxmpp/config"
	"fmt"
)

func (a *ActionTemplate) Pubsub_DoDisco() bool {

	a.user.DoRespond(a.conn,
		fmt.Sprintf("<iq type=\"result\" xmlns=\"jabber:client\" from=\"pubsub.%s\" id=\"%s\" to=\"%s\">"+
			"<query xmlns=\"http://jabber.org/protocol/disco#info\">"+
			"<identity type=\"service\" category=\"pubsub\" name=\"PubSub acs-clustered\" />"+
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
			"<x type=\"result\" xmlns=\"jabber:x:data\">"+
			"<field type=\"hidden\" var=\"FORM_TYPE\">"+
			"<value>http://jabber.org/network/serverinfo</value>"+
			"</field>"+
			"<field type=\"text-multi\" var=\"abuse-addresses\">"+
			a.GetAbuseInfo()+
			"</field>"+
			"</x>"+
			"</query>"+
			"</iq>",
			config.Config.Server.Domain,
			a.data.Id,
			a.user.FullAddr,
		), a.data.Id)
	return true
}

func (a *ActionTemplate) ActionPubsub() bool {
	if a.data.Pubsub.Items.Node == "storage:bookmarks" {
		return a.ActionGetBookmarks()
	}
	return true
}

func (a *ActionTemplate) ActionGetBookmarks() bool {
	//todo implement roster notes XEP-0145
	//log.Println("DO Storage Bookmarks")

	a.user.DoRespond(a.conn,
		a.GetResultHeader()+"<pubsub xmlns='http://jabber.org/protocol/pubsub'>"+
			"<items node='storage:bookmarks'/>"+
			"</pubsub>"+
			"</iq>", a.data.Id)
	return true
}
