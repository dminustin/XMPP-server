package actions

import (
	"amfxmpp/modules"
	"crypto/tls"
	"log"
)

func DoRespond(conn *tls.Conn, msg string, id string) error {
	if id != "" {
		modules.WriteQueChan(id, msg)
	}
	i, err := conn.Write([]byte(msg))
	if false {
		log.Printf("wrote %s bytes for %s", i, msg)
	}
	return err
}
