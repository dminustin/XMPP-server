package actions

func (a *ActionTemplate) ActionRosterGetList() bool {

	DoRespond(a.conn,
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

func (a *ActionTemplate) ActionRoster() bool {
	if a.data.Query.Roster.XMLName.Space == "roster:delimiter" {
		DoRespond(a.conn,
			a.GetResultHeader()+"</iq>", a.data.Id)
		return true
	}
	DoRespond(a.conn,
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
