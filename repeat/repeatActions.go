package repeat

import (
	"github.com/empijei/wapty/ui"
	"github.com/empijei/wapty/ui/apis"
)

var uiRepeater ui.Subscription

func init() {
	uiRepeater = ui.Subscribe(apis.REPEATCHANNEL)
}
