package actions

import (
	"fmt"
)

func (a *ActionTemplate) ActionBind() bool {
	a.user.Resource = a.data.Bind.Payload.Content
	if a.user.Resource != "" {
		a.user.FullAddr = a.user.UID + "/" + a.user.Resource
	} else {
		a.user.FullAddr = a.user.UID
	}
	DoRespond(a.conn,
		fmt.Sprintf("<iq id='%s' type='result'>"+
			"<bind xmlns='urn:ietf:params:xml:ns:xmpp-bind'>"+
			"<jid>%s</jid>"+
			"</bind>"+
			"</iq>",
			a.data.Id, a.user.FullAddr))
	return true
}
