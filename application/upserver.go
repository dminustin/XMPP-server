package application

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"io"
	//"encoding/pem"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	//"fmt"
	"strings"

	conf "amfxmpp/config"
	//"crypto/tls"
	"amfxmpp/modules"
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
	http.HandleFunc(conf.Config.FileServer.PutPath, UploadServerHandler)
	http.HandleFunc(conf.Config.FileServer.DownloadPath, DownloadServerHandler)

	err = server.ListenAndServeTLS(conf.Config.Server.Public_key, conf.Config.Server.Private_key)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

func DownloadServerHandler(w http.ResponseWriter, req *http.Request) {
	//log.Println("RECV")
	//log.Println("HEADER", req.Header)
	//log.Println("METHOD", req.Method)
	//log.Println("METHOD", req.URL)

	defer req.Body.Close()
	_, _ = ioutil.ReadAll(req.Body)

	spl := strings.Split(req.URL.Path, "/")
	key := spl[len(spl)-1]
	uid := spl[len(spl)-1]
	uid = uid[:2]
	var re = regexp.MustCompile(`([^a-z0-9\\._])`)
	key = re.ReplaceAllString(key, ``)
	uid = re.ReplaceAllString(uid, ``)

	if (len(key) < 10) || (len(uid) < 1) {
		http.Error(w, "File not found.", 404)
		return
	}

	dir := conf.Config.FileServer.FileStoragePath + uid + "/" + key
	fName := dir
	e := filepath.Ext(dir)
	log.Println(e, dir)
	if len(e) != 4 {
		matches, err := filepath.Glob(dir + ".*")
		log.Println(dir + ".*")
		if (err != nil) || (len(matches) < 1) {
			http.Error(w, "File not found.", 404)
			return
		}

		fName = matches[0]
	}
	log.Println(fName)
	Openfile, err := os.Open(fName)
	defer Openfile.Close() //Close after function return
	if err != nil {
		//File not found, send 404
		http.Error(w, "File not found.", 404)
		return
	}

	//File is found, create and send the correct headers

	//Get the Content-Type of the file
	//Create a buffer to store the header of the file in
	FileHeader := make([]byte, 512)
	//Copy the headers into the FileHeader buffer
	Openfile.Read(FileHeader)
	//Get content type of file
	FileContentType := http.DetectContentType(FileHeader)
	fmt.Println(FileContentType)

	//Get the file size
	FileStat, _ := Openfile.Stat()                     //Get info from file
	FileSize := strconv.FormatInt(FileStat.Size(), 10) //Get file size as a string

	//Send the headers
	//w.Header().Set("Content-Disposition", "attachment; filename="+key)
	w.Header().Set("Content-Type", FileContentType)
	w.Header().Set("Content-Length", FileSize)

	//Send the file
	//We read 512 bytes from the file already, so we reset the offset back to 0
	Openfile.Seek(0, 0)
	io.Copy(w, Openfile)

}

func UploadServerHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("SENT")
	log.Println("HEADER", req.Header)
	//log.Println("METHOD", req.Method)
	//log.Println("METHOD", req.URL)

	defer req.Body.Close()

	if req.ContentLength > 5*1024*1024 {
		http.Error(w, `Payload Too Large`, 413)
		return
	}

	contentType := req.Header.Get(`Content-Type`)
	log.Println("Trying to upload", contentType)

	fileExt := ``

	val, ok := conf.Config.AllowedTypes[contentType]
	if ok {
		fileExt = `.` + val
	} else {
		http.Error(w, `Unsupported Media Type`, 415)
		return
	}

	file, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Unauthorized access", 403)
		return
	}
	//log.Println(string(b))
	spl := strings.Split(req.URL.Path, "/")
	key := spl[len(spl)-1]
	uid := spl[len(spl)-2]

	var re = regexp.MustCompile(`([^a-z0-9\\._])`)
	key = re.ReplaceAllString(key, ``)
	uid = re.ReplaceAllString(uid, ``)

	var dbUpl dbUploads
	err = modules.DB.Get(&dbUpl, `select * from xmpp_uploads where hash=? and from_id=? limit 1`, key, uid)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	dir := string(key)[0:2]

	dir = conf.Config.FileServer.FileStoragePath + dir + "/"

	filename := key + fileExt

	err = ioutil.WriteFile(dir+filename, file, 0644)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	} else {

	}

	modules.DB.Exec(`update xmpp_uploads set filename=? where hash=? and from_id=? limit 1`, filename, key, uid)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Okey!\n"))

}
