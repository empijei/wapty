package intercept

import (
	"bytes"
	"testing"

	"github.com/empijei/wapty/ui/apis"
)

/*
functions
   -historyLoop()
   -newReqResp(rawReq []byte) : uint
*/

func genDummyStatus() *History {
	var st History
	metadata := &apis.ReqRespMetaData{
		ID:          0,
		Host:        "host",
		Method:      "method",
		Path:        "path",
		Params:      true,
		Edited:      true,
		Status:      "status",
		Length:      42,
		ContentType: "content type",
		Extension:   "extension",
		TLS:         true,
		IP:          "ip",
		Port:        "port",
		Cookies:     "cookie",
		Time:        "time",
	}

	rr := []*ReqResp{&ReqResp{
		ID:           0,
		MetaData:     metadata,
		RawReq:       []byte("raw request"),
		RawRes:       []byte("ras response"),
		RawEditedReq: []byte("raw edited request"),
		RawEditedRes: []byte("ras edited response"),
	}}

	st.Count = 1
	st.ReqResps = rr

	return &st
}

var expected = []byte(`{"Count":1,"ReqResps":[{"ID":0,"MetaData":{"ID":0,"Host":"host","Method":"method","Path":"path","Params":true,"Edited":true,"Status":"status","Length":42,"ContentType":"content type","Extension":"extension","TLS":true,"IP":"ip","Port":"port","Cookies":"cookie","Time":"time"},"RawReq":"cmF3IHJlcXVlc3Q=","RawRes":"cmFzIHJlc3BvbnNl","RawEditedReq":"cmF3IGVkaXRlZCByZXF1ZXN0","RawEditedRes":"cmFzIGVkaXRlZCByZXNwb25zZQ=="}]}
`)

func TestSave(t *testing.T) {
	st := genDummyStatus()

	b := bytes.NewBuffer(nil)
	err := st.Save(b)
	if err != nil {
		t.Error(err)
	}

	if bytes.Compare(b.Bytes(), expected) != 0 {
		t.Errorf("History save failed: expected \n<%s>\nbut got \n<%s>", string(expected), string(b.Bytes()))
	}
}

func TestLoad(t *testing.T) {
	in := bytes.NewBuffer(expected)
	var st History

	err := st.Load(in)
	if err != nil {
		t.Error(err)
	}

	out := bytes.NewBuffer(nil)
	err = st.Save(out)
	if err != nil {
		t.Error(err)
	}

	if bytes.Compare(out.Bytes(), expected) != 0 {
		t.Errorf("History load failed: expected \n<%s>\nbut got \n<%s>", string(expected), string(out.Bytes()))
	}

}
