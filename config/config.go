package config

import (
	"gopkg.in/ini.v1"
	"os"
)

type ConfigStruct struct {
	Server struct {
		Domain      string
		Public_key  string
		Private_key string
		Port        int
		Ip          string
	}

	Mysql struct {
		Host     string
		Database string
		Login    string
		Password string
	}

	Password struct {
		Salt1 string
		Salt2 string
	}
}

var Config ConfigStruct

func Init() {
	cfg, err := ini.Load("./app.ini")
	if err != nil {
		println(err)
		os.Exit(1)
	}
	Config.Server.Domain = cfg.Section("server").Key("domain").String()
	Config.Server.Public_key = cfg.Section("server").Key("public_key").String()
	Config.Server.Domain = cfg.Section("server").Key("domain").String()
	Config.Server.Private_key = cfg.Section("server").Key("private_key").String()
	Config.Server.Port, _ = cfg.Section("server").Key("port").Int()
	Config.Server.Ip = cfg.Section("server").Key("ip").String()

	Config.Mysql.Host = cfg.Section("mysql").Key("host").String()
	Config.Mysql.Database = cfg.Section("mysql").Key("database").String()
	Config.Mysql.Login = cfg.Section("mysql").Key("login").String()
	Config.Mysql.Password = cfg.Section("mysql").Key("password").String()

	Config.Password.Salt1 = cfg.Section("password").Key("salt1").String()
	Config.Password.Salt2 = cfg.Section("password").Key("salt2").String()

}
