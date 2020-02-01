package application

import (
	"crypto/tls"
	"regexp"

	//"encoding/xml"
	"amfxmpp/actions"
	appconfig "amfxmpp/config"
	"amfxmpp/modules"
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
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
	var cert, _ = tls.LoadX509KeyPair(appconfig.Config.Server.Public_key, appconfig.Config.Server.Private_key)
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

	listener, err := net.Listen("tcp", ":"+fmt.Sprintf("%v", appconfig.Config.Server.Port))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer listener.Close()

	// Handle each connection.
	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		go TCPAnswer(conn)
	}
}

func TCPAnswer(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 512)
	for {
		//log.Print("server: conn: waiting")
		n, err := conn.Read(buf)
		if err != nil {
			if err != nil {
				log.Printf("server: conn: read: %s", err)
			}
			break
		}
		var s = string(buf[:n])
		//log.Printf("server: conn: echo %q\n", s)

		mtype := ParseMSGType(s)
		//log.Printf("Command [%s]", mtype)
		switch mtype {
		case "stream":
			{
				//		log.Println("Starting stream")
				WriteMessage(actions.MessageHelloReply(), conn)
				break
			}
		case "starttls":
			{
				//		log.Printf("Starting TLS")
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
	_, err := conn.Write(msg)
	if err != nil {
		log.Println("[ERROR] " + fmt.Sprintf("%s", err))
	}
	//log.Printf("server: conn: wrote %d bytes", n)

}

func doAction(msgType string, s string, conn *tls.Conn, user *modules.User) (result, fatal bool) {
	switch msgType {
	case "auth":
		{
			res := actions.ActionAuth(s, conn, user)
			//log.Printf("[%s]: %s", msgType, res)
			return res, !res
			break
		}
	case "stream":
		{
			log.Println(s)
			if user.Authorized {
				user.DoRespond(conn, actions.MessageAfterLogged(), "")
			}
			break
		}
	case "iq":
		{
			if user.Authorized {
				//log.Printf("Starting IQ")
				actions.ActionIQ(s, conn, user)

			}
			break
		}
	case "presence":
		{
			if user.Authorized {
				//log.Printf("Starting Presence")
				actions.ActionPresence(s, conn, user)

			}
			break
		}
	case "message":
		{
			//log.Printf("Starting Messaging")
			if user.Authorized {
				actions.ActionMessage(s, conn, user)

			}
			break
		}
	default:
		{
			if s != "" {
				log.Printf("[%s] unknown cmd %s", msgType, s)
				//user.DoRespond(conn, "", "")
			}
		}
	}
	return true, false
}

func handleTLSConnection(unenc_conn net.Conn) {
	user := &modules.User{
		Authorized:           false,
		LastServerRequest:    0,
		PayLoad:              "",
		Resource:             "",
		FullAddr:             "",
		ReadyForInteractions: false,
		LastMessageID:        "0",
		LastSentMessageID:    0,
	}
	WriteMessage(actions.MessageProceedTLS(), unenc_conn)

	//log.Printf("%s", "Start server")
	conn := tls.Server(unenc_conn, &config)
	//log.Printf("%s", "Start handshake")

	_ = conn.Handshake()
	//fmt.Println(err)
	//log.Printf("%s", "End handshake")
	_, _ = conn.Write(actions.MessageHelloReply())
	//log.Printf("server: conn: wrote %d bytes", n)
	var buffer = make([]byte, 16384)

	connEstablished := false

	for {
		bytesRead, err := conn.Read(buffer)
		if err != nil {
			conn.Close()
			return
		}
		var in_string = string(buffer[:bytesRead])
		//log.Printf("server: conn: echo %q\n", in_string)

		tags := []string{
			"message",
			"iq",
			"presence",
		}

		parsed := false
		for _, tag := range tags {
			x := fmt.Sprintf(`(<%s[ |>].*?</%s>)`, tag, tag)
			re := regexp.MustCompile(x)
			results := re.FindAllString(in_string, 9999999)
			for _, s := range results {
				parsed = true
				mtype := ParseMSGType(s)
				//log.Printf("Command [%s]", mtype)
				doAction(mtype, s, conn, user)
			}
		}
		if !parsed {
			mtype := ParseMSGType(in_string)
			doAction(mtype, in_string, conn, user)
		}

		//log.Println("AUTH:", user.Authorized, "READY:", user.ReadyForInteractions, user)
		if (user.Authorized) && (user.ReadyForInteractions) {
			if !connEstablished {
				//log.Println("Start Interactions routine")
				connEstablished = true
				go modules.DoServerInteractions(user, conn)
			}

			//PrintMemUsage()

		}
	}
}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
