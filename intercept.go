package main

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"path"

	"github.com/empijei/WAPTy/mitm"
)

var (
	hostname, _ = os.Hostname()

	dir      = path.Join(os.Getenv("HOME"), ".mitm")
	keyFile  = path.Join(dir, "ca-key.pem")
	certFile = path.Join(dir, "ca-cert.pem")
)

type MayBeRequest struct {
	Req *http.Request
	Err error
}
type PendingRequest struct {
	OriginalRequest *http.Request
	ModifiedRequest chan *MayBeRequest
}

var RequestQueue chan *PendingRequest

type MayBeResponse struct {
	Res *http.Response
	Err error
}
type PendingResponse struct {
	OriginalResponse *http.Response
	OriginalRequest  *http.Request
	ModifiedResponse chan *MayBeResponse
}

var ResponseQueue chan *PendingResponse

var Done chan struct{}

func DispatchLoop() {
	for {
		select {
		case preq := <-RequestQueue:
			r := preq.OriginalRequest
			req, err := httputil.DumpRequest(r, true)
			if err != nil {
				preq.ModifiedRequest <- &MayBeRequest{Err: err}
				break
			}
			//fmt.Printf("%s", req)
			_ = ioutil.WriteFile("tmp.request", req, 0644)
			_, _ = stdin.ReadString('\n')
			log.Println("Continued")
			//TODO chech this error
			editedRequestFile, _ := os.Open("tmp.request")
			//TODO chech this error
			editedRequest, _ := http.ReadRequest(bufio.NewReader(editedRequestFile))
			preq.ModifiedRequest <- &MayBeRequest{Req: editedRequest}

		case presp := <-ResponseQueue:
			res := presp.OriginalResponse
			res.ContentLength = -1
			res.Header.Set("Content-Length", "-1")
			rawRes, err := httputil.DumpResponse(res, true)
			if err != nil {
				log.Println("Error while dumping response" + err.Error())
				presp.ModifiedResponse <- &MayBeResponse{Err: err}
				break
			}
			_ = ioutil.WriteFile("tmp.response", rawRes, 0644)
			_, _ = stdin.ReadString('\n')
			log.Println("Continued")
			//TODO chech this error
			editedResponseFile, _ := os.Open("tmp.response")
			editedResponseBuffer := bufio.NewReader(editedResponseFile)
			//TODO chech this error
			editedResponse, _ := http.ReadResponse(editedResponseBuffer, presp.OriginalRequest)
			//TODO adjust content length Header?
			//tmp, _ := httputil.DumpResponse(editedResponse, true)
			//fmt.Printf("%s", tmp)
			presp.ModifiedResponse <- &MayBeResponse{Res: editedResponse, Err: err}

		case <-Done:
			return
		}
	}
}

type ResponseInterceptor struct {
	WrappedRT http.RoundTripper
}

func (ri *ResponseInterceptor) RoundTrip(req *http.Request) (res *http.Response, err error) {
	res, err = ri.WrappedRT.RoundTrip(req)
	log.Println("Response intercepted")
	if err != nil {
		return
	}
	ModifiedResponse := make(chan *MayBeResponse)
	ResponseQueue <- &PendingResponse{ModifiedResponse: ModifiedResponse, OriginalRequest: req, OriginalResponse: res}
	mayBeRes := <-ModifiedResponse
	res, err = mayBeRes.Res, mayBeRes.Err
	return
}

var stdin *bufio.ReadWriter

func init() {
	RequestQueue = make(chan *PendingRequest)
	ResponseQueue = make(chan *PendingResponse)
	Done = make(chan struct{})
	stdin = bufio.NewReadWriter(bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdin))
}
func main() {
	go DispatchLoop()
	ca, err := loadCA()
	if err != nil {
		log.Fatal(err)
	}
	modifiedTransport := ResponseInterceptor{WrappedRT: http.DefaultTransport}
	p := &mitm.Proxy{
		CA: &ca,
		TLSServerConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
			//CipherSuites: []uint16{tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA},
		},
		Wrap:      intercept,
		Transport: &modifiedTransport,
	}
	log.Fatal(http.ListenAndServe(":8080", p))
}

func loadCA() (cert tls.Certificate, err error) {
	// TODO(kr): check file permissions
	cert, err = tls.LoadX509KeyPair(certFile, keyFile)
	if os.IsNotExist(err) {
		cert, err = genCA()
	}
	if err == nil {
		cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
	}
	return
}

func genCA() (cert tls.Certificate, err error) {
	err = os.MkdirAll(dir, 0700)
	if err != nil {
		return
	}
	certPEM, keyPEM, err := mitm.GenerateCA(hostname)
	if err != nil {
		return
	}
	cert, _ = tls.X509KeyPair(certPEM, keyPEM)
	err = ioutil.WriteFile(certFile, certPEM, 0400)
	if err == nil {
		err = ioutil.WriteFile(keyFile, keyPEM, 0400)
	}
	return cert, err
}

func intercept(upstream http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Request intercepted")
		ModifiedRequest := make(chan *MayBeRequest)
		RequestQueue <- &PendingRequest{OriginalRequest: r, ModifiedRequest: ModifiedRequest}
		mayBeReq := <-ModifiedRequest
		if mayBeReq.Err != nil {
			upstream.ServeHTTP(w, r)
			return
		}
		upstream.ServeHTTP(w, mayBeReq.Req)
	})
}
