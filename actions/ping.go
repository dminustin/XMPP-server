package actions

import (
	"fmt"
)

func (a *ActionTemplate) ActionPong() bool {
	DoRespond(a.conn,
		fmt.Sprintf("<iq from='%s' to='%s' id='%s' type='result'/>",
			a.data.To, a.user.FullAddr, a.data.Id), a.data.Id)
	return true
}
