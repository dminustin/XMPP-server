package application

import (
	"log"
	"net"
	"os"
)

//Upload server

func UploadServer(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 512)
	for {
		log.Print("server: conn: waiting")
		n, err := conn.Read(buf)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		log.Println("\n\n\n\n\n")
		log.Println(n, buf)
		log.Println("\n\n\n\n\n")
		os.Exit(1)
	}
}
