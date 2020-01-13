package modules

import (
	"amfxmpp/config"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"strconv"
	"time"
)

type User struct {
	ID         string
	UID        string
	Authorized bool
	Resource   string
	FullAddr   string
}

type userStruct struct {
	ID    int    `db:"id"`
	Name  string `db:"nickname"`
	State string `db:"user_state"`
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

func (u *User) TryToAuth(login, password string) (bool, string) {
	passw := hashPassw(password)
	t, err := strconv.ParseInt(login, 10, 32)
	if err != nil {
		return false, "Invalid ID"
	}

	var id = fmt.Sprintf("%v", t)

	var user userStruct
	err = DB.Get(&user, "select id, user_state, nickname from users where id=? and user_password=?", id, passw)
	if err != nil {
		return false, "Unknown user"
	}
	log.Printf("%s", user)
	if user.State != "active" {
		return false, "Your state is " + user.State
	}
	return true, "You are welcome"
}
