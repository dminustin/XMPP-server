package main

import (
	app "amfxmpp/application"
	config "amfxmpp/config"
	modules "amfxmpp/modules"
)

func main() {
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
	app.Init()

}
