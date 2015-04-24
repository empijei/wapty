package mitm

import (
	"crypto/tls"
	"log"
	"net/http"
	"testing"
)

func ExampleProxy(t *testing.T) {
	ca, err := loadCA()
	if err != nil {
		log.Fatal(err)
	}

	p := &Proxy{
		CA: &ca,
		Wrap: func(upstream http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				log.Println("Got Content-Type:", r.Header.Get("Content-Type"))
				upstream.ServeHTTP(w, r)
			})
		},
	}
	listenAndServe(p)
}

func loadCA() (cert tls.Certificate, err error) { panic("example only") }
func listenAndServe(_ http.Handler)             { panic("example only") }
