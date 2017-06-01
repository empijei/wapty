package intercept

import (
	"github.com/empijei/wapty/ui"
	"github.com/empijei/wapty/ui/apis"
)

var uiHistory *ui.Subscription

func init() {
	uiHistory = ui.Subscribe(apis.HISTORYCHANNEL)
}
