package modules

import (
	"amfxmpp/config"
	"crypto/md5"
	"crypto/tls"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type User struct {
	ID                   string
	UID                  string
	Authorized           bool
	ReadyForInteractions bool
	Resource             string
	FullAddr             string
	LastServerRequest    int64
	PayLoad              string
	LastMessageID        string
}

type userStruct struct {
	ID            int            `db:"id"`
	Name          string         `db:"nickname"`
	State         string         `db:"user_state"`
	Lastlogin     int64          `db:"last_login"`
	LastMessageID sql.NullString `db:"message_id"`
}

func getMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func hashPassw(passw string) string {
	m1 := getMD5Hash(passw + config.Config.Password.Salt1)
	m2 := getMD5Hash(m1 + config.Config.Password.Salt2)
	return m2
}

func (u *User) ChangeResource(res string) {
	if !u.ReadyForInteractions {
		u.FullAddr = res
		tmp := strings.Split(res, "/")
		u.Resource = tmp[1]
		u.ReadyForInteractions = true
	}
}

func (u *User) GetUploadToken() string {
	/**
	This token will correct next hour
	*/
	t := strconv.FormatInt(time.Now().Unix(), 10)

	m1 := getMD5Hash("Upload" + t + u.ID + config.Config.Password.Salt1 + fmt.Sprintf("%i", rand.Int()))
	m2 := getMD5Hash(m1 + config.Config.Password.Salt2 + fmt.Sprintf("%i", rand.Int()))
	m3 := getMD5Hash(m2 + config.Config.Password.Salt1 + fmt.Sprintf("%i", rand.Int()))
	return m3 + "." + t
}

func (u *User) TryToAuth(login, password string, resource string) (bool, string) {
	passw := hashPassw(password)
	t, err := strconv.ParseInt(login, 10, 32)
	if err != nil {
		return false, "Invalid ID"
	}

	var id = fmt.Sprintf("%v", t)
	resource = strconv.Quote(strings.ToLower(resource))

	var user userStruct
	err = DB.Get(&user, `select 
		(SELECT max(messages.id) FROM messages where messages.from_user=? or messages.to_user=?) as message_id,
		users.id, users.user_state, users.nickname, 
		UNIX_TIMESTAMP(COALESCE(xmpp_sessions.last_login, xmpp_sessions.last_login, 0)) as last_login 
		from users 
		left join messages on messages.to_user=users.id
		left join xmpp_sessions on xmpp_sessions.user_id=users.id and 
		xmpp_sessions.user_resource=? 
		where users.id=? and user_password=?`,
		id, id,
		resource, id, passw)
	if err != nil {
		log.Println(err)
		return false, "Unknown user"
	}

	if user.State != "active" {
		return false, "Your state is " + user.State
	}
	u.LastServerRequest = 0

	if user.LastMessageID.Valid {
		u.LastMessageID = user.LastMessageID.String
	} else {
		u.LastMessageID = "0"
	}

	return true, "You are welcome"
}

func (u *User) DoRespond(conn *tls.Conn, msg string, id string) error {
	if id != "" {
		WriteQueChan(id, msg)
	}

	if u.PayLoad != "" {
		msg = msg + u.PayLoad
		log.Println("\n\n\n[PAYLOAD]\n")
		log.Println(u.PayLoad)

		u.PayLoad = ""
	}
	if msg == "" {
		return nil
	}
	i, err := conn.Write([]byte(msg))
	if false {
		log.Printf("wrote %s bytes for %s", i, msg)
	}

	return err
}
