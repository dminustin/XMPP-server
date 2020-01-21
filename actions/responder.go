package actions

//import (
//	"amfxmpp/modules"
//	"crypto/tls"
//	"log"
//)

//func (u *modules.User) DoRespond(conn *tls.Conn, msg string, id string) error {
//	if id != "" {
//		modules.WriteQueChan(id, msg)
//	}
//
//	if u.PayLoad != "" {
//		msg = msg + u.PayLoad
//		u.PayLoad = ""
//	}
//	i, err := conn.Write([]byte(msg))
//
//	if false {
//		log.Printf("wrote %s bytes for %s", i, msg)
//	}
//
//	return err
//}
