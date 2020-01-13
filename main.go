package main

import (
	app "amfxmpp/application"
	config "amfxmpp/config"
)

func main() {
	config.Init()
	app.Init()

}
