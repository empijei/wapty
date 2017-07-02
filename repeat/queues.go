package repeat

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"

	"github.com/empijei/wapty/ui/apis"
)

var done = make(chan struct{})

// RepeaterLoop is the main loop for the repeater. It listens on apis.REPEATCHANNEL
// for calls and executes them
func RepeaterLoop() {
	for {
		select {
		case cmd := <-uiRepeater.RecChannel():
			switch cmd.Action {
			case apis.CREATE:
				r := NewRepeater()
				status.Add(r)
			case apis.GO:
				uiRepeater.Send(handleGo(&cmd))
			case apis.GET:
				uiRepeater.Send(handleGet(&cmd))
			default:
				log.Println("Unknown repeater action: " + cmd.Action)
			}
		case <-done:
			return
		}
	}
}

func handleGo(cmd *apis.Command) *apis.Command {
	var host string
	var tls bool
	var ri int
	err := cmd.UnpackArgs(
		[]apis.ArgName{apis.ENDPOINT, apis.TLS, apis.ID},
		&host, &tls, &ri,
	)
	if err != nil {
		log.Println(err)
		return apis.Err(err)
	}
	body := bytes.NewBuffer(cmd.Payload)
	status.RLock()
	defer status.RUnlock()
	if len(status.Repeats) <= ri || ri < 0 {
		err := "Repeater out of range"
		log.Println(err)
		return apis.Err(err)
	}
	r := status.Repeats[ri]
	var res io.Reader
	if res, err = r.repeat(body, host, tls); err != nil {
		log.Println(err)
		return apis.Err(err)
	}
	resbuf, err := ioutil.ReadAll(res)
	if err != nil {
		return apis.Err(err)
	}
	cmd.Payload = resbuf
	return cmd
}

func handleGet(cmd *apis.Command) *apis.Command {
	var ri, itemn int
	err := cmd.UnpackArgs(
		[]apis.ArgName{apis.ID, apis.SUBID},
		&ri, &itemn,
	)
	status.RLock()
	defer status.RUnlock()
	if len(status.Repeats) <= ri {
		log.Println("Repeater out of range")
		return apis.Err(err)
	}
	r := status.Repeats[ri]
	if len(r.History) <= itemn {
		err := "Repeater item out of range"
		log.Println(err)
		return apis.Err(err)
	}
	repitem, err := json.Marshal(r.History[itemn])
	if err != nil {
		err := "Error while marshaling repeat item"
		log.Println(err)
		return apis.Err(err)
	}
	cmd.Payload = repitem
	return cmd
}
