package webroot

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

//TODO add mimetypes
var webFiles = map[string]string{}

//This is pretty clever, thanks creack https://stackoverflow.com/a/21596576
func LoadRoutes() {
	for fileName, content := range webFiles {
		contentCpy := content
		log.Printf("Loading /" + fileName)
		var contentType string
		if strings.HasSuffix(fileName, ".js") {
			contentType = "application/javascript"
			log.Println(" as " + contentType)
		} else {
			if strings.HasSuffix(fileName, ".css") {
				contentType = "text/css"
				log.Println(" as " + contentType)
			}
			log.Println()
		}
		http.HandleFunc("/"+fileName, func(w http.ResponseWriter, r *http.Request) {
			if contentType != "" {
				w.Header().Set("Content-Type", contentType)
			}
			fmt.Fprintf(w, "%s\n", contentCpy)
		})
	}
}
