package modules

import (
	"amfxmpp/config"
	"amfxmpp/structs"
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

type MessageStruct struct {
	ID         string `db:"id"`
	Message    string `db:"message"`
	Attachment string `db:"attachment"`
	AttID      string `db:"att_id"`
	FromUser   string `db:"from_user"`
	ToUser     string `db:"to_user"`
}

func (u *User) GetFriendsUpdates() []structs.DBRosterStruct {
	s := `SELECT 
` + config.Config.Tables.TableFriendship + `.*,
users.nickname
FROM ` + config.Config.Tables.TableFriendship + `
JOIN users ON users.id = ` + config.Config.Tables.TableFriendship + `.friend_id
AND ` + config.Config.Tables.TableFriendship + `.state in ("friend", "memory")
JOIN xmpp_sessions
ON xmpp_sessions.user_id=` + config.Config.Tables.TableFriendship + `.user_id
AND xmpp_sessions.last_login<` + config.Config.Tables.TableFriendship + `.contact_state_date
/*AND xmpp_sessions.user_resource=?*/
WHERE ` + config.Config.Tables.TableFriendship + `.user_id=?`
	var statuses []structs.DBRosterStruct
	err := DB.Select(&statuses, s, u.ID) //	u.Resource

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
		message,
		IF(messages_attachements.message_id IS NOT NULL, CONCAT(" https://`+
			config.Config.Server.Domain+":"+config.Config.FileServer.DownloadPort+
			config.Config.FileServer.DownloadPath+
			`", SUBSTRING(messages_attachements.filename,1,2), "/",
                                                        messages_attachements.filename
                                                       ), "") as attachment,
		COALESCE(messages_attachements.filename, "") as att_id,
		messages.from_user, messages.to_user 
		FROM messages 
		LEFT JOIN messages_attachements
		ON 
		messages_attachements.message_id=messages.id AND
		messages_attachements.to_id = messages.to_user AND
		messages_attachements.from_id = messages.from_user
		 
		 
		 where
		messages.id>%s and
		(messages.to_user=%s /* or messages.from_user=%s */)
		order by messages.id asc
		`,
		u.LastMessageID,
		u.ID,
		u.ID,
	)

	err := DB.Select(&messages, sq)

	if err != nil {
		log.Println(err)
	}

	//srch := "https://" + config.Config.Server.Domain+":"+config.Config.FileServer.DownloadPort+
	//	config.Config.FileServer.DownloadPath
	//re := regexp.MustCompile(`(^(Uploaded file )`+srch+`)`)
	lmid, _ := strconv.ParseInt(u.LastMessageID, 10, 64)
	for _, msg := range messages {
		z, _ := strconv.ParseInt(msg.ID, 10, 64)
		if z > lmid {
			u.LastMessageID = msg.ID
		}

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
			set last_login=NOW(), user_id=?, user_resource=? , last_msg_read_id=?
			on duplicate key update 
			last_login=NOW(),
			last_msg_read_id=?
			`,
			//u.LastServerRequest,
			u.ID,
			u.Resource, u.LastMessageID, u.LastMessageID,
			//u.LastServerRequest,
		)
		//log.Println(u.FullAddr, u.LastMessageID)
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
	attach := ""
	attid := ``
	if message.Attachment != "" {

		attach = `<x xmlns="jabber:x:oob">
		<url>` + message.Attachment + `</url>
		</x>`
		attid = "-" + message.AttID
	}
	if message.Message == "Uploaded file" {
		message.Message = message.Attachment
	}

	isCarbon := false
	canCarbon := false

	if message.FromUser == user.ID {
		e, _ := strconv.ParseInt(message.ID, 10, 64)
		isCarbon = true
		if e > user.LastSentMessageID {
			canCarbon = true
		}
	}

	msgID := "srvmsg-" + message.ID + attid

	m := fmt.Sprintf(`<message from="%s" xmlns="jabber:client" id="%s" to="%s" type="chat">
<body>%s</body>
<origin-id xmlns="urn:xmpp:sid:0" id="%s" />
<active xmlns="http://jabber.org/protocol/chatstates" />
<request xmlns="urn:xmpp:receipts" />`+attach+`
</message>`,
		message.FromUser+"@"+config.Config.Server.Domain,
		msgID,
		user.FullAddr,
		message.Message,
		msgID,
	)
//do nothing for a while
canCarbon = false;
	if isCarbon && canCarbon {

		m = `<message xmlns='jabber:client' 
			from='` + user.ID + "@" + config.Config.Server.Domain + `' to='` + user.FullAddr + `'
         type='chat'>
  <received xmlns='urn:xmpp:carbons:2'>
<forwarded xmlns='urn:xmpp:forward:0'>
` + m + `
</forwarded>
  </received>
</message>
`

	}

	_, err := conn.Write([]byte(m))
	return err

}
