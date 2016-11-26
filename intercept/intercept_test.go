package intercept

import (
	"encoding/json"
	"net/http"
	"os"
	"testing"
)

func setup() {
}

func shutdown() {}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func jsonEqual(a interface{}, b interface{}) (equal bool, as string, bs string) {
	buf, err := json.MarshalIndent(a, " ", " ")
	if err != nil {
		panic(err)
	}
	as = string(buf)
	buf, err = json.MarshalIndent(b, " ", " ")
	if err != nil {
		panic(err)
	}
	bs = string(buf)
	equal = as == bs
	return
}

func reqEqual(a *http.Request, b *http.Request) bool {
	//TODO IMPLEMENT THIS but only on exported fields
	//return reflect.DeepEqual(a, b)
	return true
}
