//go:generate go run $GOROOT/src/crypto/tls/generate_cert.go -host "example.com,127.0.0.1" -ca -ecdsa-curve P256
//go:generate sh -c "go-bindata -o cert_test.go -pkg mitm *.pem"
package mitm

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

var hostname, _ = os.Hostname()

var (
	nettest = flag.Bool("nettest", false, "run tests over network")
)

func init() {
	flag.Parse()
}

var (
	caCert = MustAsset("cert.pem")
	caKey  = MustAsset("key.pem")
)

func testProxy(t *testing.T, setupReq func(req *http.Request), wrap func(http.Handler) http.Handler, downstream http.HandlerFunc, checkResp func(*http.Response)) {
	ds := httptest.NewTLSServer(downstream)
	defer ds.Close()

	rootCAs := x509.NewCertPool()
	if !rootCAs.AppendCertsFromPEM(caCert) {
		panic("can't add cert")
	}

	ca, err := tls.X509KeyPair(caCert, caKey)
	if err != nil {
		panic(err)
	}
	ca.Leaf, err = x509.ParseCertificate(ca.Certificate[0])
	if err != nil {
		panic(err)
	}
	log.Printf("ca.Leaf.IPAddresses: %v", ca.Leaf.IPAddresses)

	p := &Proxy{
		CA: &ca,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		TLSServerConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
		Wrap: wrap,
	}

	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatal("Listen:", err)
	}
	defer l.Close()

	go func() {
		if err := http.Serve(l, p); err != nil {
			if !strings.Contains(err.Error(), "use of closed network") {
				t.Fatal("Serve:", err)
			}
		}
	}()

	t.Logf("requesting %q", ds.URL)
	req, err := http.NewRequest("GET", ds.URL, nil)
	if err != nil {
		t.Fatal("NewRequest:", err)
	}
	setupReq(req)

	c := &http.Client{
		Transport: &http.Transport{
			Proxy: func(r *http.Request) (*url.URL, error) {
				u := *r.URL
				u.Scheme = "https"
				u.Host = l.Addr().String()
				return &u, nil
			},
			TLSClientConfig: &tls.Config{
				RootCAs: rootCAs,
			},
		},
	}

	resp, err := c.Do(req)
	if err != nil {
		t.Fatal("Do:", err)
	}
	checkResp(resp)
}

func Test(t *testing.T) {
	const xHops = "X-Hops"

	testProxy(t, func(req *http.Request) {
		// req.Host = "example.com"
		req.Header.Set(xHops, "a")
	}, func(upstream http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			hops := r.Header.Get("X-Hops") + "b"
			r.Header.Set("X-Hops", hops)
			upstream.ServeHTTP(w, r)
		})
	}, func(w http.ResponseWriter, r *http.Request) {
		hops := r.Header.Get(xHops) + "c"
		w.Header().Set(xHops, hops)
	}, func(resp *http.Response) {
		const w = "abc"
		if g := resp.Header.Get(xHops); g != w {
			t.Errorf("want %s to be %s, got %s", xHops, w, g)
		}
	})
}

func TestNet(t *testing.T) {
	if !*nettest {
		t.Skip()
	}

	var wrapped bool
	testProxy(t, func(req *http.Request) {
		nreq, _ := http.NewRequest("GET", "https://mitmtest.herokuapp.com/", nil)
		*req = *nreq
	}, func(upstream http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			wrapped = true
			upstream.ServeHTTP(w, r)
		})
	}, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("this shouldn't be hit")
	}, func(resp *http.Response) {
		if !wrapped {
			t.Errorf("expected wrap")
		}
		got, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal("ReadAll:", err)
		}
		if code := resp.StatusCode; code != 200 {
			t.Errorf("want code 200, got %d", code)
		}
		if g := string(got); g != "ok\n" {
			t.Errorf("want ok, got %q", g)
		}
	})
}
