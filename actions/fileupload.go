package actions

import (
	"amfxmpp/config"
	"amfxmpp/modules"
	"fmt"
	"log"
	"os"
)

func (a *ActionTemplate) FileUpload_DoDisco() bool {
	//todo implement Upload Max File Size into config
	a.user.DoRespond(a.conn,
		fmt.Sprintf(
			"<iq type=\"result\" xmlns=\"jabber:client\" from=\"upload.%s\" id=\"%s\" to=\"%s\">"+
				"<query xmlns=\"http://jabber.org/protocol/disco#info\">"+
				"<identity type=\"file\" category=\"store\" name=\"HTTP File Upload component\" />"+
				"<feature var=\"urn:xmpp:http:upload:0\" />"+
				"<x type=\"result\" xmlns=\"jabber:x:data\">"+
				"<field type=\"hidden\" var=\"FORM_TYPE\">"+
				"<value>http://jabber.org/network/serverinfo</value>"+
				"</field>"+
				a.GetAbuseInfo()+
				"</x>"+
				"<x type=\"result\" xmlns=\"jabber:x:data\">"+
				"<field var=\"FORM_TYPE\">"+
				"<value>urn:xmpp:http:upload:0</value>"+
				"</field>"+
				"<field var=\"max-file-size\">"+
				"<value>%s</value>"+
				"</field>"+
				"</x>"+
				"</query>"+
				"</iq>",
			config.Config.Server.Domain,
			a.data.Id,
			a.user.FullAddr,
			"52428800", //5 Mb

		), a.data.Id)

	return true
}

func (a *ActionTemplate) ActionRequestFileUpload() bool {
	log.Printf("[Do uploading] %s", a.data)
	uplHash := a.user.GetUploadToken()

	_, err := modules.DB.Exec(`insert into xmpp_uploads
	set hash=?, from_id=?`, uplHash, a.user.ID)

	if err != nil {
		log.Println(err)
		os.Exit(1)
		a.user.DoRespond(a.conn,
			fmt.Sprintf(`<result type="error" to="%s" id="%s" />`, a.user.FullAddr, a.data.Id),
			a.data.Id)
		return true
	}

	//todo put "PUT URL" && "GET URL" into config
	a.user.DoRespond(a.conn,
		fmt.Sprintf("<iq to=\"%s\" xmlns=\"jabber:client\" id=\"%s\" from=\"upload.%s\" type=\"result\">"+
			"<slot xmlns=\"urn:xmpp:http:upload:0\">"+
			"<put url=\"%s\" />"+
			"<get url=\"https://download/"+uplHash+"\" />"+
			"</slot>"+
			"</iq><a xmlns=\"urn:xmpp:sm:3\" h=\"25\" />",
			a.user.FullAddr, a.data.Id, config.Config.Server.Domain,
			"https://"+config.Config.Server.Domain+":"+config.Config.FileServer.UploadPort+config.Config.FileServer.PutPath+a.user.ID+"/"+uplHash,
		), a.data.Id)
	return true
}
