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
