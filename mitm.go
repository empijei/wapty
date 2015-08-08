package mitm

import (
	"crypto/tls"
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"sync"
	"time"
)

type ServerParam struct {
	CA        *tls.Certificate // the Root CA for generatng on the fly MITM certificates
	TLSConfig *tls.Config      // a template TLS config for the server.
}

// A ServerConn is a net.Conn that holds its clients SNI header in ServerName
// after the handshake.
type ServerConn struct {
	*tls.Conn

	// ServerName is set during Conn's handshake to the client's requested
	// server name set in the SNI header. It is not safe to access across
	// multiple goroutines while Conn is performing the handshake.
	ServerName string
}

// Server wraps cn with a ServerConn configured with p so that during its
// Handshake, it will generate a new certificate using p.CA. After a successful
// Handshake, its ServerName field will be set to the clients requested
// ServerName in the SNI header.
func Server(cn net.Conn, p ServerParam) *ServerConn {
	conf := new(tls.Config)
	if p.TLSConfig != nil {
		*conf = *p.TLSConfig
	}
	sc := new(ServerConn)
	conf.GetCertificate = func(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
		sc.ServerName = hello.ServerName
		return getCert(p.CA, hello.ServerName)
	}
	sc.Conn = tls.Server(cn, conf)
	return sc
}

type listener struct {
	net.Listener
	ca   *tls.Certificate
	conf *tls.Config
}

// NewListener returns a net.Listener that generates a new cert from ca for
// each new Accept. It uses SNI to generate the cert, and herefore only
// works with clients that send SNI headers.
//
// This is useful for building transparent MITM proxies.
func NewListener(inner net.Listener, ca *tls.Certificate, conf *tls.Config) net.Listener {
	return &listener{inner, ca, conf}
}

func (l *listener) Accept() (net.Conn, error) {
	cn, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}
	sc := Server(cn, ServerParam{
		CA:        l.ca,
		TLSConfig: l.conf,
	})
	return sc, nil
}

// Proxy is a forward proxy that substitutes its own certificate
// for incoming TLS connections in place of the upstream server's
// certificate.
type Proxy struct {
	// Wrap specifies a function for optionally wrapping upstream for
	// inspecting the decrypted HTTP request and response.
	Wrap func(upstream http.Handler) http.Handler

	// CA specifies the root CA for generating leaf certs for each incoming
	// TLS request.
	CA *tls.Certificate

	// TLSServerConfig specifies the tls.Config to use when generating leaf
	// cert using CA.
	TLSServerConfig *tls.Config

	// TLSClientConfig specifies the tls.Config to use when establishing
	// an upstream connection for proxying.
	TLSClientConfig *tls.Config

	// FlushInterval specifies the flush interval
	// to flush to the client while copying the
	// response body.
	// If zero, no periodic flushing is done.
	FlushInterval time.Duration

	// Director is function which modifies the request into a new
	// request to be sent using Transport. See the documentation for
	// httputil.ReverseProxy for more details. For mitm proxies, the
	// director defaults to HTTPDirector, but for transparent TLS
	// proxies it should be set to HTTPSDirector.
	Director func(*http.Request)
}

var (
	okHeader           = "HTTP/1.1 200 OK\r\n\r\n"
	noUpstreamHeader   = "HTTP/1.1 503 No Upstream\r\n\r\n"
	noDownstreamHeader = "HTTP/1.1 503 No Downstream\r\n\r\n"
	errHeader          = "HTTP/1.1 500 Internal Server Error\r\n\r\n"
)

func (p *Proxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != "CONNECT" {
		rp := &httputil.ReverseProxy{
			Director:      p.Director,
			FlushInterval: p.FlushInterval,
		}
		if rp.Director == nil {
			rp.Director = HTTPDirector
		}
		p.Wrap(rp).ServeHTTP(w, req)
		return
	}

	cn, _, err := w.(http.Hijacker).Hijack()
	if err != nil {
		log.Println("Hijack:", err)
		http.Error(w, "No Upstream", 503)
		return
	}
	defer cn.Close()

	_, err = io.WriteString(cn, okHeader)
	if err != nil {
		log.Println("Write:", err)
		return
	}

	sc, ok := cn.(*ServerConn)
	if !ok {
		name := dnsName(req.Host)
		if name == "" {
			log.Println("cannot determine cert name for " + req.Host)
			io.WriteString(cn, noDownstreamHeader)
			return
		}
		sc = Server(cn, ServerParam{
			CA:        p.CA,
			TLSConfig: p.TLSServerConfig,
		})
		if err := sc.Handshake(); err != nil {
			log.Println("Server Handshake:", err)
			return
		}
	}

	cc, err := p.tlsDial(req.Host, sc.ServerName)
	if err != nil {
		log.Println("tlsDial:", err)
		io.WriteString(cn, noUpstreamHeader)
		return
	}
	p.proxyMITM(sc, cc)
}

func (p *Proxy) tlsDial(addr, serverName string) (net.Conn, error) {
	conf := new(tls.Config)
	if p.TLSClientConfig != nil {
		*conf = *p.TLSClientConfig
	}
	conf.ServerName = serverName
	return tls.Dial("tcp", addr, conf)
}

func (p *Proxy) proxyMITM(upstream, downstream net.Conn) {
	var mu sync.Mutex
	dial := func(network, addr string) (net.Conn, error) {
		mu.Lock()
		defer mu.Unlock()
		if downstream == nil {
			return nil, io.EOF
		}
		cn := downstream
		downstream = nil
		return cn, nil
	}
	rp := &httputil.ReverseProxy{
		Director:      HTTPSDirector,
		Transport:     &http.Transport{DialTLS: dial},
		FlushInterval: p.FlushInterval,
	}
	ch := make(chan struct{})
	wc := &onCloseConn{upstream, func() { ch <- struct{}{} }}
	http.Serve(&oneShotListener{wc}, p.Wrap(rp))
	<-ch
}

// HTTPDirector is director designed for use in Proxy for http
// proxies.
func HTTPDirector(r *http.Request) {
	r.URL.Host = r.Host
	r.URL.Scheme = "http"
}

// HTTPSDirector is a director designed for use in Proxy for
// transparent TLS proxies.
func HTTPSDirector(req *http.Request) {
	req.URL.Host = req.Host
	req.URL.Scheme = "https"
}

// A oneShotListener implements net.Listener whos Accept only returns a
// net.Conn as specified by c followed by an error for each subsequent Accept.
type oneShotListener struct {
	c net.Conn
}

func (l *oneShotListener) Accept() (net.Conn, error) {
	if l.c == nil {
		return nil, errors.New("closed")
	}
	c := l.c
	l.c = nil
	return c, nil
}

func (l *oneShotListener) Close() error {
	return nil
}

func (l *oneShotListener) Addr() net.Addr {
	return l.c.LocalAddr()
}

// A onCloseConn implements net.Conn and calls its f on Close.
type onCloseConn struct {
	net.Conn
	f func()
}

func (c *onCloseConn) Close() error {
	if c.f != nil {
		c.f()
		c.f = nil
	}
	return c.Conn.Close()
}

// dnsName returns the DNS name in addr, if any.
func dnsName(addr string) string {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return ""
	}
	return host
}

// Certificates are cached locally to avoid unnecessary regeneration
const certCacheMaxSize = 1000

var (
	certCache      = make(map[*tls.Certificate]map[string]*tls.Certificate)
	certCacheMutex sync.RWMutex
)

func getCert(ca *tls.Certificate, host string) (*tls.Certificate, error) {
	if c := getCachedCert(ca, host); c != nil {
		return c, nil
	}
	cert, err := GenerateCert(ca, host)
	if err != nil {
		return nil, err
	}
	cacheCert(ca, host, cert)
	return cert, nil
}

func getCachedCert(ca *tls.Certificate, host string) *tls.Certificate {
	certCacheMutex.RLock()
	defer certCacheMutex.RUnlock()

	if certCache[ca] == nil {
		return nil
	}
	cert := certCache[ca][host]
	if cert == nil || cert.Leaf.NotAfter.Before(time.Now()) {
		return nil
	} else {
		return cert
	}
}

func cacheCert(ca *tls.Certificate, host string, cert *tls.Certificate) {
	certCacheMutex.Lock()
	defer certCacheMutex.Unlock()

	if certCache[ca] == nil || len(certCache[ca]) > certCacheMaxSize {
		certCache[ca] = make(map[string]*tls.Certificate)
	}
	certCache[ca][host] = cert
}
