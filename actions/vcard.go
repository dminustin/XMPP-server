package actions

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"amfxmpp/config"
	"amfxmpp/modules"
)

func (a *ActionTemplate) ActionVCardGet() bool {

	tmp := strings.Split(a.data.To, `@`)
	if len(tmp) < 2 {
		return false
	}

	i, err := strconv.ParseInt(tmp[0], 10, 64)
	if err != nil {
		log.Println(err)
		return false
	}
	log.Println(tmp)
	info, err := modules.GetUserByID(fmt.Sprintf("%v", i))
	if err != nil {
		log.Println(err)
		return false
	}

	photo := ``
	if info.Photo.BinVal != `` {
		photo = `<PHOTO>
        <BINVAL>` + info.Photo.BinVal + `</BINVAL>
        <TYPE>` + info.Photo.Type + `</TYPE>
        </PHOTO>`
	}

	result := `<vCard xmlns="vcard-temp">
<NICKNAME>` + info.Nickname + `</NICKNAME>
<BDAY>1976-02-03</BDAY>
<JABBERID>` + fmt.Sprintf("%v", i) + `@` + config.Config.Server.Domain + `</JABBERID>
<URL>https://` + config.Config.Server.Domain + `/members/` + fmt.Sprintf("%v", i) + `</URL>
<DESC>` + info.AboutMe + `</DESC>
` + photo + `
</vCard>`

	a.user.DoRespond(a.conn,
		fmt.Sprintf("<iq from='%s' to='%s' id='%s' type='result'>"+result+"</iq>",
			a.data.To, a.user.FullAddr, a.data.Id), a.data.Id)

	return true
}

func (a *ActionTemplate) ActionVCardSet() bool {
	log.Println("Do not forget to implement SET VCard method!!")
	a.user.DoRespond(a.conn,
		fmt.Sprintf("<iq from='%s' to='%s' id='%s' type='result'/>",
			a.data.To, a.user.FullAddr, a.data.Id), a.data.Id)
	return true
}
