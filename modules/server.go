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
	ToUser   string `db:"to_user"`
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
/*AND xmpp_sessions.user_resource=?*/
WHERE xmpp_roster.owner_id=?`
	var statuses []structs.DBRosterStruct
	err := DB.Select(&statuses, s, u.ID)//	u.Resource

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return statuses
}

func (u *User) GetNewMessages() []MessageStruct {
	var messages []MessageStruct
	//todo get messages from and to user

	sq := fmt.Sprintf(
		`
SELECT messages.id, 
IF(messages_attachements.message_id IS NOT NULL, CONCAT(messages.message, " https://`+
			config.Config.Server.Domain+
			config.Config.FileServer.DownloadPath+
			`", SUBSTRING(messages_attachements.filename,1,2), "/",
                                                        messages_attachements.filename
                                                       ), messages.message) as message,

messages.from_user, messages.to_user 
FROM messages 
LEFT JOIN messages_attachements
ON 
messages_attachements.message_id=messages.id AND
messages_attachements.to_id = messages.to_user AND
messages_attachements.from_id = messages.from_user
		 join xmpp_sessions
		 on xmpp_sessions.user_id = messages.to_user
		 and xmpp_sessions.user_resource = "%s"
		 and 
		 (
			(%s !='0' AND messages.id> %s)
		 /*OR (messages.date_create>=xmpp_sessions.last_login)*/
			)
		 where messages.to_user=%s
		order by messages.id asc
		`,
		u.Resource,
		u.LastMessageID, u.LastMessageID,
		u.ID,
	)

	err := DB.Select(&messages, sq)

	if err != nil {
		fmt.Println(err)
	}
	return messages

}

func DoServerInteractions(u *User, conn *tls.Conn) {
	for {
		time.Sleep(time.Second * 10)

		messages := u.GetNewMessages()
		for _, msg := range messages {
			err := ActionPullMessage(&msg, conn, u)
			if err != nil {
				return
			}
			u.LastMessageID = msg.ID
		}

		updates := u.GetFriendsUpdates()
		for _, msg := range updates {
			err := ActionPullFriendUpdate(&msg, conn, u)
			if err != nil {
				return
			}
		}

		u.LastServerRequest = time.Now().Unix()
		DB.Exec(`insert into xmpp_sessions 
			set last_login=NOW(), user_id=?, user_resource=? 
			on duplicate key update last_login=NOW()`,
			//u.LastServerRequest,
			u.ID,
			u.Resource,
			//u.LastServerRequest,
		)

	}
}

func ActionPullFriendUpdate(message *structs.DBRosterStruct, conn *tls.Conn, user *User) error {
	msgID := "srvstatus-" + fmt.Sprintf("%v", time.Now().Unix()) + "-" + message.UserID

	state := ""
	if message.ContactState.String != "active" {
		state = "<show>" + message.ContactState.String + "</show>"
	} else {
		message.ContactStatusMessage.String = ""
	}

	status := ""
	if len(message.ContactStatusMessage.String) > 0 {
		status = "<status >" + message.ContactStatusMessage.String + "</status >"
	}

	s := fmt.Sprintf(`<presence from="%s@%s" xmlns="jabber:client" id="%s" to="%s">
<priority>50</priority>
%s
%s
<c node="http://gajim.org" hash="sha-1" xmlns="http://jabber.org/protocol/caps" ver="0oVRDLJYyCnbS13aGaP3gSFUU/o=" />
</presence>
`,
		message.UserID, config.Config.Server.Domain, msgID, user.FullAddr,
		state, status,
		message.ContactStatusMessage.String,
	)

	s = s + `<message from="` + message.UserID + "@" + config.Config.Server.Domain + `" 
	xmlns="jabber:client" id="a` + msgID + `" to="` + user.FullAddr + `">
<event xmlns="http://jabber.org/protocol/pubsub#event">
<items node="http://jabber.org/protocol/activity">
<item id="current">
<activity xmlns="http://jabber.org/protocol/activity" />
</item>
</items>
</event>
</message>
<message from="` + message.UserID + "@" + config.Config.Server.Domain + `" 
	xmlns="jabber:client" id="b` + msgID + `" to="` + user.FullAddr + `">
<event xmlns="http://jabber.org/protocol/pubsub#event">
<items node="http://jabber.org/protocol/mood">
<item id="current">
<mood xmlns="http://jabber.org/protocol/mood" />
</item>
</items>
</event>
</message>`

	//user.PayLoad = user.PayLoad + "\r\n" + s
	_, err := conn.Write([]byte(s))
	return err
}

func ActionPullMessage(message *MessageStruct, conn *tls.Conn, user *User) error {
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
	_, err := conn.Write([]byte(m))
	return err

}
