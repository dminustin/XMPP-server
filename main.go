package main

import (
	"log"
	"time"

	app "amfxmpp/application"
	"amfxmpp/config"
	"amfxmpp/modules"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("App started")
	config.Init()
	modules.InitDB(
		config.Config.Mysql.Login,
		config.Config.Mysql.Password,
		config.Config.Mysql.Host,
		config.Config.Mysql.Port,
		config.Config.Mysql.Database,
	)
	var m modules.User
	m.GetUploadToken()
	go app.InitUploadServer()
	go app.Init()
	for {
		time.Sleep(time.Second * 60)
	}

}
