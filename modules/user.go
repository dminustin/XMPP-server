package modules

import (
	"amfxmpp/config"
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"log"
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
}

type userStruct struct {
	ID        int    `db:"id"`
	Name      string `db:"nickname"`
	State     string `db:"user_state"`
	Lastlogin int64  `db:"last_login"`
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

func (u *User) GetUploadToken() string {
	/**
	This token will correct next hour
	*/
	t := strconv.FormatInt(time.Now().Unix(), 10)

	m1 := getMD5Hash("Upload" + t + u.ID + config.Config.Password.Salt1)
	m2 := getMD5Hash(m1 + config.Config.Password.Salt2)

	return m2 + "." + t
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
	err = DB.Get(&user, "select "+
		"users.id, users.user_state, users.nickname,"+
		" UNIX_TIMESTAMP(COALESCE(xmpp_sessions.last_login, xmpp_sessions.last_login, 0)) as last_login "+
		" from users"+
		" left join xmpp_sessions on xmpp_sessions.user_id=users.id and "+
		" xmpp_sessions.user_resource=?"+
		" where id=? and user_password=?", resource, id, passw)
	if err != nil {
		return false, "Unknown user"
	}

	if user.State != "active" {
		return false, "Your state is " + user.State
	}
	u.LastServerRequest = 0
	return true, "You are welcome"
}

func (u *User) DoRespond(conn *tls.Conn, msg string, id string) error {
	if id != "" {
		WriteQueChan(id, msg)
	}

	if u.PayLoad != "" {
		msg = msg + u.PayLoad
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
