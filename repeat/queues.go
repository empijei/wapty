package repeat

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"strconv"

	"github.com/empijei/wapty/cli/lg"
	"github.com/empijei/wapty/ui/apis"
)

var done = make(chan struct{})

// RepeaterLoop is the main loop for the repeater. It listens on apis.CHN_REPEAT
// for calls and executes them
func RepeaterLoop() {
	for {
		select {
		case cmd := <-uiRepeater.RecChannel():
			switch cmd.Action {
			case apis.RPT_CREATE:
				uiRepeater.Send(handleCreate(&cmd))
			case apis.RPT_GO:
				uiRepeater.Send(handleGo(&cmd))
			case apis.RPT_GET:
				uiRepeater.Send(handleGet(&cmd))
			default:
				lg.Infof("Unknown repeater action: %s", cmd.Action)
			}
		case <-done:
			return
		}
	}
}

func handleCreate(cmd *apis.Command) *apis.Command {
	r := NewRepeater()
	id := status.Add(r)
	cmd.Args = map[apis.ArgName]string{apis.ARG_ID: strconv.Itoa(id)}

	//TODO reply with repeater ID
	return cmd
}

func handleGo(cmd *apis.Command) *apis.Command {
	var host string
	var tls bool
	var ri int
	err := cmd.UnpackArgs(
		[]apis.ArgName{apis.ARG_ENDPOINT, apis.ARG_TLS, apis.ARG_ID},
		&host, &tls, &ri,
	)
	if err != nil {
		lg.Error(err)
		return apis.Err(err)
	}
	body := bytes.NewBuffer(cmd.Payload)
	status.RLock()
	defer status.RUnlock()
	if len(status.Repeats) <= ri || ri < 0 {
		err := "Repeater out of range"
		lg.Error(err)
		return apis.Err(err)
	}
	r := status.Repeats[ri]
	var res io.Reader
	var id int
	if res, id, err = r.repeat(body, host, tls); err != nil {
		lg.Error("%s\n", err.Error())
		return apis.Err(err)
	}
	resbuf, err := ioutil.ReadAll(res)
	if err != nil {
		return apis.Err(err)
	}
	cmd.Payload = resbuf
	cmd.Args[apis.ARG_SUBID] = strconv.Itoa(id)
	return cmd
}

func handleGet(cmd *apis.Command) *apis.Command {
	var ri, itemn int
	err := cmd.UnpackArgs(
		[]apis.ArgName{apis.ARG_ID, apis.ARG_SUBID},
		&ri, &itemn,
	)
	status.RLock()
	defer status.RUnlock()
	if len(status.Repeats) <= ri {
		lg.Infof("Repeater out of range")
		return apis.Err(err)
	}
	r := status.Repeats[ri]
	if len(r.History) <= itemn {
		err := "Repeater item out of range"
		lg.Error(err)
		return apis.Err(err)
	}
	repitem, err := json.Marshal(r.History[itemn])
	if err != nil {
		err := "Error while marshaling repeat item"
		lg.Error(err)
		return apis.Err(err)
	}
	cmd.Payload = repitem
	return cmd
}
