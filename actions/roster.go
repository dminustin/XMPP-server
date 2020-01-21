package actions

import (
	"amfxmpp/config"
	"amfxmpp/modules"
	"amfxmpp/structs"

	"fmt"
	"os"
)

func (a *ActionTemplate) ActionRosterGetList() bool {

	var sel []structs.DBRosterStruct
	err := modules.DB.Select(&sel, "select xmpp_roster.owner_id,"+
		"xmpp_roster.user_id , "+
		"xmpp_roster.relation , "+
		"xmpp_roster.contact_state , "+
		"xmpp_roster.contact_state_date , "+
		"xmpp_roster.contact_state_date , "+
		"xmpp_roster.contact_status_message , "+
		" users.nickname from xmpp_roster "+
		"join users on users.id=xmpp_roster.user_id "+
		"where xmpp_roster.owner_id=? and xmpp_roster.relation='friend'", a.user.ID)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	result := a.GetResultHeader() +
		"<query xmlns=\"jabber:iq:roster\" ver=\"" + a.GetDateVersion() + "\">"

	fmt.Println(sel)

	for t, r := range sel {
		fmt.Println(t, r)
		result = result + fmt.Sprintf(`<item jid="%s@%s" subscription="both" name="%s" />`,
			r.UserID,
			config.Config.Server.Domain,
			r.Nickname,
		)
	}

	a.user.DoRespond(a.conn,
		result+"</query></iq>",
		a.data.Id)
	a.user.ReadyForInteractions = true
	return true
}

func (a *ActionTemplate) ActionRoster() bool {
	if a.data.Query.Roster.XMLName.Space == "roster:delimiter" {
		a.user.DoRespond(a.conn,
			a.GetResultHeader()+"</iq>", a.data.Id)
		return true
	}
	return a.ActionRosterGetList()
}
