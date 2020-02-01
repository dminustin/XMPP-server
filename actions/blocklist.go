package actions

import "log"

func (a *ActionTemplate) ActionBlockList() bool {
	if a.data.Blocklist.Xmlns == "urn:xmpp:blocking" && a.data.Type == "get" {
		return a.ActionGetBlockList()
	}
	if a.data.Blocklist.Xmlns == "urn:xmpp:blocking" && a.data.Type == "set" {
		return a.ActionSetBlockList()
	}
	a.user.DoRespond(a.conn,
		a.GetResultHeader()+
			"<query xmlns=\"jabber:iq:roster\" ver=\"8edf7fb47032e2acb95908651113d861\">"+
			"<item jid=\"2@tematicon.club\" subscription=\"both\" name=\"222\" />"+
			"<item jid=\"3@tematicon.club\" subscription=\"both\" name=\"333\" />"+
			"<item jid=\"4@tematicon.club\" subscription=\"both\" name=\"444\" />"+
			"</query>"+
			"</iq>",
		a.data.Id)
	return true
}

func (a *ActionTemplate) ActionGetBlockList() bool {
	a.user.DoRespond(a.conn,
		"<iq type='result' id='"+a.data.Id+"'>"+
			"<blocklist xmlns='urn:xmpp:blocking'>"+
			"<item jid='romeo@montague.net'/>"+
			"<item jid='iago@shakespeare.lit'/>"+
			"</blocklist>"+
			"</iq>", a.data.Id)
	return true

}

func (a *ActionTemplate) ActionSetBlockList() bool {

	blockJID := a.data.Blocklist.Item.Content
	log.Println("Set block to " + blockJID)
	a.user.DoRespond(a.conn,
		"<iq type='result' id='"+a.data.Id+"' />", a.data.Id)
	return true

}
