package application

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	//"fmt"
	"strings"

	//"encoding/pem"
	"io/ioutil"

	conf "amfxmpp/config"
	//"crypto/tls"
	"amfxmpp/modules"
	"log"
	"net/http"
)

type dbUploads struct {
	Hash     string         `db:"hash"`
	FromID   string         `db:"from_id"`
	Filename sql.NullString `db:"filename"`
	Regdate  string         `db:"regdate"`
}

func InitUploadServer() {
	var cert, _ = tls.LoadX509KeyPair(conf.Config.Server.Public_key, conf.Config.Server.Private_key)
	caCert, err := ioutil.ReadFile(conf.Config.Server.Public_key)
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	config := &tls.Config{
		MinVersion:   tls.VersionTLS10,
		Certificates: []tls.Certificate{cert},
		CipherSuites: []uint16{tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA},
		ServerName:               conf.Config.Server.Domain,
		PreferServerCipherSuites: true,
		InsecureSkipVerify:       true,
		ClientCAs:                caCertPool,
		ClientAuth:               tls.NoClientCert,
		//ClientAuth: tls.RequireAndVerifyClientCert,
	}

	//config.BuildNameToCertificate()
	//modules.InitDB(
	//	conf.Config.Mysql.Login,
	//	conf.Config.Mysql.Password,
	//	conf.Config.Mysql.Host,
	//	conf.Config.Mysql.Port,
	//	conf.Config.Mysql.Database,
	//)
	server := &http.Server{
		TLSConfig: config,
		Addr:      ":7778", //+ conf.Config.FileServer.UploadPort,
	}

	server.Handler = http.DefaultServeMux
	http.HandleFunc(conf.Config.FileServer.PutPath, HelloServer)

	err = server.ListenAndServeTLS(conf.Config.Server.Public_key, conf.Config.Server.Private_key)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

func HelloServer(w http.ResponseWriter, req *http.Request) {
	log.Println("SENT")

	log.Println("HEADER", req.Header)

	log.Println("METHOD", req.Method)

	log.Println("METHOD", req.URL)

	defer req.Body.Close()

	file, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	//log.Println(string(b))
	spl := strings.Split(req.URL.Path, "/")
	key := spl[len(spl)-1]
	uid := spl[len(spl)-2]

	var dbUpl dbUploads
	err = modules.DB.Get(&dbUpl, `select * from xmpp_uploads where hash=? and from_id=? limit 1`, key, uid)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	dir := string(key)[0:2]

	dir = conf.Config.FileServer.FileStoragePath + dir + "/"

	filename := key + ".jpg"

	err = ioutil.WriteFile(dir+filename, file, 0644)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	modules.DB.Exec(`update xmpp_uploads set filename=? where hash=? and from_id=? limit 1`, filename, key, uid)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Okey!\n"))

}
