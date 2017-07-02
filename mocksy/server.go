package mocksy

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func mocksyHandler(rw http.ResponseWriter, req *http.Request) {
	resp := FindMatching(req)
	//if err != nil {
	//rw.WriteHeader(http.StatusInternalServerError)
	//log.Println(err)
	//fmt.Fprintln(rw, "mocksy: internal server error :(")
	//return
	//}
	rw.WriteHeader(http.StatusOK)
	fmt.Fprintln(rw, string(resp.Value))
}

func StartServer(port string) error {
	if err := loadResponseHistory(); err != nil {
		return fmt.Errorf("mocksy: error starting server:\n\t%s", err.Error())
	}
	http.HandleFunc("/", mocksyHandler)
	http.ListenAndServe(port, nil)
	return nil
}

func loadResponseHistory() error {
	// Import data (TODO: should not be hardcoded here, also should check errors)
	fname := "test.xml"
	file, err := os.Open(fname)
	if err != nil {
		return fmt.Errorf("File not found: %s", fname)
	}

	// Go refuses to parse any XML whose version is != "1.0". Burp sometimes
	// declares XML 1.1, albeit it uses no 1.1-only features, so we trick
	// the XML parser into parsing our "invalid" XML by skipping the XML header.
	buf := make([]byte, len(`<?xml version="1.x"?>`))
	if n, err := file.Read(buf); err == nil && n == len(buf) {
	} else {
		return fmt.Errorf("mocksy: error importing data: header skip failed.\n")
	}

	items, err := BurpImport(file)
	if err != nil {
		return fmt.Errorf("mocksy: error importing data:\n\t%s", err.Error())
	}

	for _, item := range items.Items {
		AddToHistory(item)
	}
	log.Printf("Loaded %d Request-Response pairs.\n", HistoryLength())

	return nil
}
