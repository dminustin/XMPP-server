package actions

import (
	"fmt"
)

func (a *ActionTemplate) ActionEnable() bool {
	DoRespond(a.conn,
		fmt.Sprintf("<enabled xmlns='%s'/>",
			a.data.Enable.Xmlns,
		))
	return true
}
