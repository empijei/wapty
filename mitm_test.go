package mitm

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
)

var hostname, _ = os.Hostname()

var (
	nettest = flag.Bool("nettest", false, "run tests over network")
)

func init() {
	flag.Parse()
}

func genCA() (cert tls.Certificate, err error) {
	certPEM, keyPEM, err := GenCA(hostname)
	if err != nil {
		return tls.Certificate{}, err
	}
	cert, err = tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return tls.Certificate{}, err
	}
	cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
	return cert, err
}

func testProxy(t *testing.T, ca *tls.Certificate, setupReq func(req *http.Request), wrap func(http.Handler) http.Handler, downstream http.HandlerFunc, checkResp func(*http.Response)) {
	ds := httptest.NewTLSServer(downstream)
	defer ds.Close()

	p := &Proxy{
		CA: ca,
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
				InsecureSkipVerify: true,
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

	ca, err := genCA()
	if err != nil {
		t.Fatal("loadCA:", err)
	}

	testProxy(t, &ca, func(req *http.Request) {
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

	ca, err := genCA()
	if err != nil {
		t.Fatal("loadCA:", err)
	}

	var wrapped bool
	testProxy(t, &ca, func(req *http.Request) {
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
