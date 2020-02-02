package actions

import (
	"fmt"
	"log"

	"amfxmpp/utils"
)

func (a *ActionTemplate) ActionBind() bool {
	isUpdatedRes := false
	if a.data.Bind.Payload.Content != a.user.Resource && a.data.Bind.Payload.Content != "" {
		a.data.Bind.Payload.Content = utils.QuoteText(a.data.Bind.Payload.Content)
		isUpdatedRes = true
		a.user.Resource = a.data.Bind.Payload.Content
		a.user.FullAddr = a.user.UID + "/" + a.user.Resource
	} else {
		a.user.FullAddr = a.user.UID
	}

	if isUpdatedRes {
		a.user.UpdateUserFromSessionTable()
	}

	log.Println("start bind", a.user.FullAddr, a.data.Bind.Payload.Content)

	a.user.DoRespond(a.conn,
		fmt.Sprintf("<iq id='%s' type='result'>"+
			"<bind xmlns='urn:ietf:params:xml:ns:xmpp-bind'>"+
			"<jid>%s</jid>"+
			"</bind>"+
			"</iq>",
			a.data.Id, a.user.FullAddr), a.data.Id)
	return true
}
