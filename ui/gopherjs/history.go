// +build js

package main

import (
	"encoding/json"
	"fmt"
	"reflect"

	. "github.com/empijei/wapty/ui/apis"
	"github.com/gopherjs/gopherjs/js"
)

var (
	tmpHistory = make(map[int]map[string]*js.Object)
)

func handleHistory(msg Command) {
	switch msg.Action {
	case HST_METADATA:
		var md ReqRespMetaData
		err := json.Unmarshal([]byte(msg.Args[HST_METADATA]), &md)
		if err != nil {
			panic(err)
		}

		if rowMap, ok := tmpHistory[md.ID]; !ok {
			row := historyTbody.Call("insertRow", -1)
			tmp := make(map[string]*js.Object)

			val := reflect.Indirect(reflect.ValueOf(md))
			for i := 0; i < val.Type().NumField(); i++ {
				typeField := val.Type().Field(i)
				cell := row.Call("insertCell", -1)
				valueField := val.Field(i).Interface()
				cell.Set("innerText", fmt.Sprintf("%v", valueField))
				tmp[typeField.Name] = cell
			}
			tmpHistory[md.ID] = tmp
		} else {
			val := reflect.Indirect(reflect.ValueOf(md))
			for i := 0; i < val.Type().NumField(); i++ {
				typeField := val.Type().Field(i)
				cell := rowMap[typeField.Name]
				valueField := val.Field(i).Interface()
				cell.Set("innerText", fmt.Sprintf("%v", valueField))
			}
			//FIXME this is commented because it looks like the page
			//receives the same metadata multiple times.
			//delete(tmpHistory, md.Id)
		}
	case HST_FETCH:
		var rr ReqResp
		err := json.Unmarshal(msg.Payload, &rr)
		if err != nil {
			panic(err)
		}
		//TODO check if is printable, otherwise show hex
		historyReqBuffer.SetTextContent(string(rr.RawReq))
		historyResBuffer.SetTextContent(string(rr.RawRes))
	}
}

//DOM Stuff
var (
	historyTbody     *js.Object
	historyReqBuffer *DomElement
	historyResBuffer *DomElement
)

func init() {
	historyTbody = js.Global.Get("historyTbody")
	historyReqBuffer = GetElementByID("historyReqBuffer")
	historyResBuffer = GetElementByID("historyResBuffer")

	hth := js.Global.Get("historyHeader")

	//This is used to make the ui adapt to backend changes in metadata
	val := reflect.Indirect(reflect.ValueOf(ReqRespMetaData{}))
	for i := 0; i < val.Type().NumField(); i++ {
		hth.Call("insertCell", -1).Set("innerText", val.Type().Field(i).Name)
	}

	js.Global.Set("hist", map[string]interface{}{
		"onHistoryCellClick": onHistoryCellClick,
	})

}

func onHistoryCellClick() {
	proxyAction(Command{
		Action:  HST_FETCH,
		Channel: HISTORYCHANNEL,
		Args:    map[ArgName]string{ARG_ID: js.Global.Get("event").Get("target").Get("parentNode").Get("childNodes").Index(0).Get("textContent").String()},
	}, true)
}
