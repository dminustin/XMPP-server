package actions

import (
	"fmt"
)

func (a *ActionTemplate) ActionEnable() bool {
	a.user.DoRespond(a.conn,
		fmt.Sprintf("<enabled xmlns='%s'/>",
			a.data.Enable.Xmlns,
		), a.data.Id)
	return true
}
