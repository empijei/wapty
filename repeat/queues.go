package repeat

import (
	"bytes"
	"io"
	"log"
	"strconv"

	"github.com/empijei/Wapty/ui"
)

func RepeaterLoop() {
	//BOOKMARK
	select {
	case cmd := <-uiRepeater.DataChannel:
		switch action := parseRepeaterAction(cmd.Action); action {
		case CREATE:
			r := NewRepeater()
			status.Add(r)
		case GO:
			handleGo(&cmd)
		case GET:
		default:
			log.Println("Unknown repeater action: " + cmd.Action)
		}
		//TODO case <-done:

	}
}

func handleGo(cmd *ui.Command) {
	if len(cmd.Args) != 3 {
		//TODO
		log.Println("Wrong number of parameters")
		return
	}
	host := cmd.Args[0]
	tls := cmd.Args[1] == "true"
	ri, err := strconv.Atoi(cmd.Args[2])
	if err != nil {
		log.Println(err)
		return
	}
	body := bytes.NewBuffer(cmd.Payload)
	status.RLock()
	defer status.RUnlock()
	if len(status.Repeats) <= ri {
		log.Println("Repeater out of range")
		return
	}
	r := status.Repeats[ri]
	var res io.Reader
	if res, err = r.Repeat(body, host, tls); err != nil {
		log.Println(err)
		return
	}
	_ = res
	//TODO send response
	//BOOKMARK
	uiRepeater.Send(
		ui.Command{},
	)
}
