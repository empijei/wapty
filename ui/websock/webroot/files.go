package webroot

import (
	"fmt"
	"log"
	"net/http"
)

var webFiles = map[string]string{}

//This is pretty clever, thanks creack https://stackoverflow.com/a/21596576
func LoadRoutes() {
	for fileName, content := range webFiles {
		log.Println("Loading /" + fileName)
		contentCpy := content
		http.HandleFunc("/"+fileName, func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "%s\n", contentCpy)
		})
	}
}
