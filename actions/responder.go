package actions

import (
	"crypto/tls"
	"log"
)

func DoRespond(conn *tls.Conn, msg string) error {
	i, err := conn.Write([]byte(msg))
	if false {
		log.Printf("wrote %s bytes for %s", i, msg)
	}
	return err
}
