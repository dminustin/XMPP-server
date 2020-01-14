package actions

import (
	"amfxmpp/config"
	"fmt"
	"log"
)

func (a *ActionTemplate) FileUpload_DoDisco() bool {
	//todo implement Upload Max File Size into config
	DoRespond(a.conn,
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
			"5242880", //5 Mb

		), a.data.Id)

	return true
}

func (a *ActionTemplate) ActionRequestFileUpload() bool {
	log.Printf("[Do uploading] %s", a.data)
	uplHash := a.user.GetUploadToken()
	//todo put "PUT URL" && "GET URL" into config
	DoRespond(a.conn,
		fmt.Sprintf("<iq to=\"%s\" xmlns=\"jabber:client\" id=\"upload.%s\" from=\"%s\" type=\"result\">"+
			"<slot xmlns=\"urn:xmpp:http:upload:0\">"+
			"<put url=\"https://%s/xmpp/upload/%s/%s\" />"+
			"<get url=\"https://%s/xmpp/files/%s\" />"+
			"</slot>"+
			"</iq>",
			a.user.FullAddr, a.data.Id, config.Config.Server.Domain,
			config.Config.Server.Domain,
			a.user.ID,
			uplHash,
			config.Config.Server.Domain,
			uplHash,
		), a.data.Id)
	return true
}
