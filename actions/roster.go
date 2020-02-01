package actions

import (
	"amfxmpp/config"
	"amfxmpp/modules"
	"amfxmpp/structs"
	"log"
	"strconv"
	"strings"

	"fmt"
	"os"
)

func (a *ActionTemplate) ActionRosterSet() bool {

	if !a.data.Query.Item.IsEmpty() {
		tmp := strings.Split(a.data.Query.Item.Jid, "@")
		id, err := strconv.ParseInt(tmp[0], 10, 64)
		if err != nil {
			return false
		}

		modules.DB.Exec(`INSERT IGNORE INTO friendship 
		set
		user_id=?,
		friend_id=?,
		state="memory",
		contact_state="away",
		contact_state_date="1970-01-01 00:00:00"
		`, a.user.ID, id,
		)

	}
	a.data.Id = a.data.Id + "-2"
	return a.ActionRosterGetList()

}

func (a *ActionTemplate) ActionRosterGetList() bool {

	var sel []structs.DBRosterStruct
	err := modules.DB.Select(&sel, `select 
		`+config.Config.Tables.TableFriendship+`.user_id,
		`+config.Config.Tables.TableFriendship+`.friend_id , 
		`+config.Config.Tables.TableFriendship+`.state , 
		`+config.Config.Tables.TableFriendship+`.contact_state , 
		`+config.Config.Tables.TableFriendship+`.contact_state_date , 
		`+config.Config.Tables.TableFriendship+`.contact_status_message , 
		 `+config.Config.Tables.TableUsers+`.nickname from `+config.Config.Tables.TableFriendship+` 
		join `+config.Config.Tables.TableUsers+` on `+config.Config.Tables.TableUsers+`.id=`+config.Config.Tables.TableFriendship+`.friend_id 
		where `+config.Config.Tables.TableFriendship+`.user_id=? 
		and `+config.Config.Tables.TableFriendship+`.state IN ('friend', 'memory')`, a.user.ID)

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	result := a.GetResultHeader() +
		`<query xmlns="jabber:iq:roster" ver="` + a.GetDateVersion() + `">`

	for _, r := range sel {

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
