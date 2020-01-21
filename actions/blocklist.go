package actions

func (a *ActionTemplate) ActionBlockList() bool {
	if a.data.Blocklist.Xmlns == "urn:xmpp:blocking" {
		return a.ActionGetBlockList()
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
