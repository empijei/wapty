package mocksy

import "net/http"

// FIXME: use a container better suited for searching. Must find an efficient key
// to do fuzzy search with requests.
type responseDB []Item

var responseHistory responseDB

func init() {
	responseHistory = make([]Item)
}

// AddToHistory inserts a pair request-response in the responseHistory.
func AddToHistory(itm Item) {
	append(responseHistory, itm)
}

func HistoryLength() int {
	return len(responseHistory)
}

// FindMatching takes an http request and returns the closest match to it
// based on the response history.
func FindMatching(req *http.Request) string {
	// TODO
	for _, item := range responseDB {
	}

	return ""
}
