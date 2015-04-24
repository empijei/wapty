package mitm

import (
	"crypto/tls"
	"log"
	"net/http"
	"testing"
)

type codeRecorder struct {
	http.ResponseWriter

	code int
}

func (w *codeRecorder) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
	w.code = code
}

func ExampleProxy(t *testing.T) {
	ca, err := loadCA()
	if err != nil {
		log.Fatal(err)
	}

	p := &Proxy{
		CA: &ca,
		Wrap: func(upstream http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				cr := &codeRecorder{ResponseWriter: w}
				log.Println("Got Content-Type:", r.Header.Get("Content-Type"))
				upstream.ServeHTTP(cr, r)
				log.Println("Got Status:", cr.code)
			})
		},
	}
	listenAndServe(p)
}

func loadCA() (cert tls.Certificate, err error) { panic("example only") }
func listenAndServe(_ http.Handler)             { panic("example only") }
