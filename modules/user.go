package modules

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"amfxmpp/config"
	"amfxmpp/utils"
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
	LastSentMessageID    int64
	UserData             UserData
}

type UserData struct {
	AvatarPath  string //path to avatar image
	AvatarHash  string
	Nickname    string
	BirthDate   string //birth date
	PersonalURL string
	Phones      struct {
		Home string
		Work string
	}
	AboutMe string
	Photo   struct {
		BinVal string
		Type   string
	}
}

type userStruct struct {
	ID            int            `db:"id"`
	Name          string         `db:"nickname"`
	State         string         `db:"user_state"`
	Lastlogin     string         `db:"last_login"`
	LastMessageID sql.NullString `db:"last_msg_read_id"`
	AvatarID      sql.NullString `db:"avatar_id"`
	BirthDate     sql.NullString `db:"bdate"`
	AboutMe       sql.NullString `db:"aboutme"`
}

func GetUserByID(id string) (UserData, error) {
	var res UserData

	log.Println("Get info for", id)

	var user userStruct
	err := DB.Get(&user, `select 
    users.id, 
    users.nickname, 
    users.avatar_id, 
    users.bdate, 
    users.aboutme 
 
    from users where id=?`, id)
	if err != nil {
		log.Println(err)
		return res, err
	}

	if user.AvatarID.Valid {
		res.Photo.Type = "image/jpeg"
		res.Photo.BinVal = getUserAvatar(user.AvatarID.String, id)
	}

	res.Nickname = user.Name

	res.BirthDate = user.BirthDate.String
	res.AboutMe = user.AboutMe.String

	return res, nil
}

func getUserAvatar(avatarID, userID string) string {
	key := utils.GetMD5Hash(avatarID + `.` + userID)
	dir := string(key)[0:2]
	dir = config.Config.FileServer.AvatarsPath + dir + "/"
	filename := key + ".jpg"
	return utils.Base64ReadFile(dir + filename)
}

func hashPassw(passw string) string {
	m1 := utils.GetMD5Hash(passw + config.Config.Password.Salt1)
	m2 := utils.GetMD5Hash(m1 + config.Config.Password.Salt2)
	return m2
}

func (u *User) ChangeResource(res string) {
	if !u.ReadyForInteractions {
		u.FullAddr = res
		tmp := strings.Split(res, "/")
		if len(tmp) > 1 {
			u.Resource = tmp[1]
			u.ReadyForInteractions = true
		}
	}
}

func (u *User) GetUploadToken() string {
	/**
	  This token will correct next hour
	*/
	t := strconv.FormatInt(time.Now().Unix(), 10)

	m1 := utils.GetMD5Hash("Upload" + t + u.ID + config.Config.Password.Salt1 + fmt.Sprintf("%i", rand.Int()))
	m2 := utils.GetMD5Hash(m1 + config.Config.Password.Salt2 + fmt.Sprintf("%i", rand.Int()))
	m3 := utils.GetMD5Hash(m2 + config.Config.Password.Salt1 + fmt.Sprintf("%i", rand.Int()))
	return m3 + "." + t
}

func (u *User) UpdateUserFromSessionTable() {

	var user userStruct
	err := DB.Get(&user, `select 
		COALESCE(xmpp_sessions.last_msg_read_id, 0) as last_msg_read_id,
		users.id, users.user_state, users.nickname, 
		UNIX_TIMESTAMP(COALESCE(xmpp_sessions.last_login, xmpp_sessions.last_login, 0)) as last_login 
		from users 
		left join messages on messages.to_user=users.id
		left join xmpp_sessions on xmpp_sessions.user_id=users.id
		and xmpp_sessions.user_resource="`+u.Resource+`" 
		where users.id="`+u.ID+`" `)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	if user.LastMessageID.Valid {
		u.LastMessageID = user.LastMessageID.String
	} else {
		u.LastMessageID = "0"
	}
	//log.Println("UPDATED USER", u)
}

func (u *User) TryToAuth(login, password string, resource string) (bool, string) {
	passw := hashPassw(password)
	t, err := strconv.ParseInt(login, 10, 32)
	if err != nil {
		return false, "Invalid ID"
	}

	var id = fmt.Sprintf("%v", t)
	//resource = strconv.Quote(strings.ToLower(resource))

	var user userStruct
	err = DB.Get(&user, `select 
		COALESCE(xmpp_sessions.last_msg_read_id, 0) as last_msg_read_id,
		users.id, users.user_state, users.nickname, users.avatar_id,
		UNIX_TIMESTAMP(COALESCE(xmpp_sessions.last_login, xmpp_sessions.last_login, 0)) as last_login 
		from users 
		left join messages on messages.to_user=users.id
		left join xmpp_sessions on xmpp_sessions.user_id=users.id
		and xmpp_sessions.user_resource="`+resource+`" 
		where users.id="`+id+`" and user_password="`+passw+`"`)

	if err != nil {
		log.Println(err)
		//os.Exit(1)
		return false, "Unknown user"
	}

	if user.AvatarID.Valid {
		u.UserData.AvatarHash = utils.Base64ToSha1(getUserAvatar(user.AvatarID.String, fmt.Sprintf("%v", user.ID)))
		//log.Println(u.UserData.AvatarHash)
	} else {
		u.UserData.AvatarHash = ``
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
		//log.Println("\n\n\n[PAYLOAD]\n")
		//log.Println(u.PayLoad)

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
