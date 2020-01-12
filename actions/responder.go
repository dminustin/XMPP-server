package actions

import "crypto/tls"

func DoRespond(conn *tls.Conn, msg string) (error) {
	_, err:= conn.Write([]byte(msg))

	return err
}
