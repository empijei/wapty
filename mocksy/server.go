package mocksy

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/empijei/wapty/cli/lg"
)

// histDir is the directory to load XML from
var histDir string = "."

// LoadResponseHistory loads all XML files found in `histDir` into the matcher's history.
// It does NOT clear the current history (use `ClearHistory()` for that).
// It does NOT recurse on the directory.
// In case of errors, it tries to load as much files as possible and reports the error afterwards.
func LoadResponseHistory(dir string) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("mocksy: Error reading history directory: %s", err.Error())
	}

	totLoaded := 0
	errorMsgs := make([]string, 0)
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".xml") {
			continue
		}
		fp, err := os.Open(file.Name())
		if err == nil {
			if err = LoadResponsesFrom(fp); err != nil {
				errorMsgs = append(errorMsgs, err.Error())
			} else {
				fmt.Fprintf(outw, "Loaded history file %s\n", file.Name())
				totLoaded++
			}
		} else {
			errorMsgs = append(errorMsgs, err.Error())
		}
	}

	fmt.Fprintf(outw, "Loaded %d files correctly (history size = %d).\n", totLoaded, HistoryLength())
	if len(errorMsgs) > 0 {
		return fmt.Errorf("mocksy: Error importing %d files: %v", len(errorMsgs), errorMsgs)
	}
	return nil
}

// LoadResponseFrom decodes an XML source and loads all req-resp pairs in the matcher's responseHistory.
func LoadResponsesFrom(source io.ReadSeeker) error {
	// Go refuses to parse any XML whose version is != "1.0". Burp sometimes
	// declares XML 1.1, albeit it uses no 1.1-only features, so we trick
	// the XML parser into parsing our "invalid" XML by skipping the XML header.
	buf := make([]byte, len(`<?xml version="1.x"?>`))
	if n, err := source.Read(buf); err == nil && n == len(buf) {
		// Check we actually skipped the XML header and, if not, rewind.
		if !bytes.Equal(buf[:len(`<?xml`)], []byte(`<?xml`)) {
			if _, err = source.Seek(0, io.SeekStart); err != nil {
				return fmt.Errorf("mocksy: error importing data: reader rewind failed.\n")
			}
		}
	} else {
		return fmt.Errorf("mocksy: error importing data: header skip failed.\n")
	}

	items, err := BurpImport(source)
	if err != nil {
		return fmt.Errorf("mocksy: error importing data:\n\t%s", err.Error())
	}

	for _, item := range items.Items {
		AddToHistory(item)
	}
	lg.Infof("Loaded %d Request-Response pairs.", len(items.Items))

	return nil
}

// SetHistDir changes the value of the history load directory. Calling `LoadResponseHistory` after this
// will load all the XML files found in the given directory.
func SetHistDir(dir string) {
	histDir = dir
}

// StartServer loads the response history from `histDir` and starts the Mocksy server on given `port`
func StartServer(port string) error {
	if err := LoadResponseHistory(histDir); err != nil {
		return fmt.Errorf("mocksy: error starting server:\n\t%s", err.Error())
	}
	http.HandleFunc("/", mocksyHandler)
	http.ListenAndServe(port, nil)
	return nil
}

// mocksyHandler is the HTTP handler of the Mocksy server
func mocksyHandler(rw http.ResponseWriter, req *http.Request) {
	resp := FindMatching(req)
	//if err != nil {
	//rw.WriteHeader(http.StatusInternalServerError)
	//lg.Error(err)
	//fmt.Fprintln(rw, "mocksy: internal server error :(")
	//return
	//}
	rw.WriteHeader(http.StatusOK)
	fmt.Fprintln(rw, string(resp.Value))
}
