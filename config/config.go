package config

import (
	"gopkg.in/ini.v1"
	"log"
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
		Port     string
	}

	Password struct {
		Salt1 string
		Salt2 string
	}

	FileServer struct {
		UploadPort      string
		DownloadPort    string
		PutPath         string
		DownloadPath    string
		FileStoragePath string
	}

	Tables struct {
		TableUsers                string
		TableMessages             string
		TableMessagesAttachements string
		TableSessions             string
		TableUploads              string
		TableFriendship           string
	}
}

var Config ConfigStruct

func Init() {
	cfg, err := ini.Load("./app.ini")
	if err != nil {
		log.Println(err)
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
	Config.Mysql.Port = cfg.Section("mysql").Key("port").String()

	Config.Password.Salt1 = cfg.Section("password").Key("salt1").String()
	Config.Password.Salt2 = cfg.Section("password").Key("salt2").String()

	Config.FileServer.UploadPort = cfg.Section("fileserver").Key("upload_port").String()
	Config.FileServer.DownloadPort = cfg.Section("fileserver").Key("download_port").String()
	Config.FileServer.PutPath = cfg.Section("fileserver").Key("put_path").String()
	Config.FileServer.DownloadPath = cfg.Section("fileserver").Key("download_path").String()
	Config.FileServer.FileStoragePath = cfg.Section("fileserver").Key("filestorage_path").String()

	Config.Tables.TableUsers = cfg.Section("tables").Key("table_users").String()
	Config.Tables.TableMessages = cfg.Section("tables").Key("table_messages").String()
	Config.Tables.TableMessagesAttachements = cfg.Section("tables").Key("table_messages_attachements").String()
	Config.Tables.TableSessions = cfg.Section("tables").Key("table_sessions").String()
	Config.Tables.TableUploads = cfg.Section("tables").Key("table_uploads").String()
	Config.Tables.TableFriendship = cfg.Section("tables").Key("table_friendship").String()

}
