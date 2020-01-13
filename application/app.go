package application

import (
	"crypto/tls"
	//"encoding/xml"
	"amfxmpp/actions"
	appconfig "amfxmpp/config"
	"amfxmpp/modules"
	"fmt"
	"log"
	"net"
	"os"
)

var (
	config tls.Config
)

type MyStruct struct {
	Name  string
	xmlns string
	Meta  map[string]interface{}
}

func Init() {
	var cert, _ = tls.LoadX509KeyPair("./.keys/fullchain3.pem", "./.keys/privkey3.pem")
	config = tls.Config{
		MinVersion:   tls.VersionTLS10,
		Certificates: []tls.Certificate{cert},
		CipherSuites: []uint16{tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA},
		ServerName: appconfig.Config.Server.Domain,
	}

	// Listen for incoming connections.
	// fmt.Sprintf(":%d", *portPtr)
	listener, err := net.Listen("tcp", appconfig.Config.Server.Ip+":"+fmt.Sprintf("%v", appconfig.Config.Server.Port))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer listener.Close()

	// Handle each connection.
	for {
		conn, err := listener.Accept()

		if err != nil {
			os.Exit(1)
		}

		go TCPAnswer(conn)
	}
}

func TCPAnswer(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 512)
	for {
		log.Print("server: conn: waiting")
		n, err := conn.Read(buf)
		if err != nil {
			if err != nil {
				log.Printf("server: conn: read: %s", err)
			}
			break
		}
		var s = string(buf[:n])
		log.Printf("server: conn: echo %q\n", s)

		mtype := ParseMSGType(s)
		log.Printf("XCommand [%s]", mtype)
		switch mtype {
		case "stream":
			{
				log.Println("Starting stream")
				WriteMessage(actions.MessageHelloReply(), conn)
				break
			}
		case "starttls":
			{
				log.Printf("Starting TLS")
				handleTLSConnection(conn)
				break
			}
		default:
			{
				log.Printf("Unknown request %s", s)
			}
		}

		if err != nil {
			log.Printf("server: error: %s", err)
			break
		}

	}

	log.Println("server: conn: closed")
}

func WriteMessage(msg []byte, conn net.Conn) {
	n, err := conn.Write(msg)
	if err != nil {
		log.Println("[ERROR] " + fmt.Sprintf("%s", err))
	}
	log.Printf("server: conn: wrote %d bytes", n)

}

func doAction(msgType string, s string, conn *tls.Conn, user *modules.User) (result, fatal bool) {
	switch msgType {
	case "auth":
		{
			res := actions.ActionAuth(s, conn, user)
			log.Printf("[%s]: %s", msgType, res)
			return res, !res
			break
		}
	case "stream":
		{
			if user.Authorized {
				actions.DoRespond(conn, actions.MessageAfterLogged())
			}
			break
		}
	case "iq":
		{
			if user.Authorized {
				log.Printf("Starting IQ")
				actions.ActionIQ(s, conn, user)

			}
			break
		}
	default:
		{
			log.Printf("[%s] unknown cmd %s", msgType, s)
		}
	}
	return true, false
}

func handleTLSConnection(unenc_conn net.Conn) {
	user := &modules.User{
		Authorized: false,
	}
	WriteMessage(actions.MessageProceedTLS(), unenc_conn)

	log.Printf("%s", "Start server")
	conn := tls.Server(unenc_conn, &config)
	log.Printf("%s", "Start handshake")

	err := conn.Handshake()
	fmt.Println(err)
	log.Printf("%s", "End handshake")
	n, _ := conn.Write(actions.MessageHelloReply())
	log.Printf("server: conn: wrote %d bytes", n)
	var buffer = make([]byte, 1024)
	for {
		bytesRead, err := conn.Read(buffer)
		if err != nil {
			conn.Close()
			return
		}
		var s = string(buffer[:bytesRead])
		log.Printf("server: conn: echo %q\n", s)
		mtype := ParseMSGType(s)
		log.Printf("Command [%s]", mtype)
		doAction(mtype, s, conn, user)
	}
}
