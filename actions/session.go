package actions

import (
	"fmt"
)

func (a *ActionTemplate) ActionSession() bool {
	DoRespond(a.conn,
		fmt.Sprintf("<iq id='%s' type='result'>"+
			"<session xmlns='urn:ietf:params:xml:ns:xmpp-session'/></iq>",
			a.data.Id))
	return true
}
