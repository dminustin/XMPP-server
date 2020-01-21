package modules

import (
	"amfxmpp/config"
	"amfxmpp/structs"
	"crypto/tls"
	"fmt"
	"os"
	"time"
)

type MessageStruct struct {
	ID       string `db:"id"`
	Message  string `db:"message"`
	FromUser string `db:"from_user"`
}

func (u *User) GetFriendsUpdates() []structs.DBRosterStruct {
	s := `SELECT 
xmpp_roster.*,
users.nickname
FROM xmpp_roster
JOIN users ON users.id = xmpp_roster.user_id
AND xmpp_roster.relation="friend"
JOIN xmpp_sessions
ON xmpp_sessions.user_id=xmpp_roster.owner_id
AND xmpp_sessions.last_login<xmpp_roster.contact_state_date
WHERE xmpp_roster.owner_id=?`
	var statuses []structs.DBRosterStruct
	err := DB.Select(&statuses, s, u.ID)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return statuses
}

func (u *User) GetNewMessages() []MessageStruct {
	var messages []MessageStruct

	err := DB.Select(&messages, "select "+
		" messages.id, messages.message, messages.from_user"+
		" from messages"+
		" join xmpp_sessions"+
		" on xmpp_sessions.user_id = messages.to_user"+
		" and xmpp_sessions.user_resource = ?"+
		" and messages.date_create>=xmpp_sessions.last_login"+
		" where messages.to_user=?", u.Resource, u.ID)

	if err != nil {
		fmt.Println(err)
	}
	return messages

}

func DoServerInteractions(u *User, conn *tls.Conn) {
	if u.Resource == "" {
		return
	}

	messages := u.GetNewMessages()
	for _, msg := range messages {
		ActionPullMessage(&msg, conn, u)
	}

	updates := u.GetFriendsUpdates()
	for _, msg := range updates {
		ActionPullFriendUpdate(&msg, conn, u)
	}

	u.LastServerRequest = time.Now().Unix()
	_, err := DB.Exec(`insert into xmpp_sessions 
			set last_login=FROM_UNIXTIME(?), user_id=?, user_resource=? 
			on duplicate key update last_login=FROM_UNIXTIME(?)`,
		u.LastServerRequest,
		u.ID,
		u.Resource,
		u.LastServerRequest,
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func ActionPullFriendUpdate(message *structs.DBRosterStruct, conn *tls.Conn, user *User) {
	msgID := "srvstatus-" + fmt.Sprintf("%v", time.Now().Unix()) + "-" + message.UserID

	s := fmt.Sprintf(`<presence from="%s@%s" xmlns="jabber:client" id="%s" to="%s@%s">
<priority>40</priority>
<show>%s</show>
<status>%s</status>
</presence>`,
		message.UserID, config.Config.Server.Domain,
		msgID,
		user.ID, config.Config.Server.Domain,
		message.ContactState.String,
		message.ContactStatusMessage.String,
	)
	user.PayLoad = user.PayLoad + "\r\n\r\n" + s
}

func ActionPullMessage(message *MessageStruct, conn *tls.Conn, user *User) {
	msgID := "srvmsg-" + message.ID
	m := fmt.Sprintf(`<message from="%s" xmlns="jabber:client" id="%s" to="%s" type="chat">
<body>%s</body>
<origin-id xmlns="urn:xmpp:sid:0" id="%s" />
<active xmlns="http://jabber.org/protocol/chatstates" />
<request xmlns="urn:xmpp:receipts" />
</message>`,
		message.FromUser+"@"+config.Config.Server.Domain,
		msgID,
		user.FullAddr,
		message.Message,
		msgID,
	)

	user.PayLoad = user.PayLoad + "\r\n\r\n" + m

}
