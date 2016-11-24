package intercept

import (
	"log"
	"net/http"
)

//Remove trailers?
//https://github.com/squid-cache/squid/blob/master/src/http/RegisteredHeadersHash.cci
var HopByHopHeaders = []string{
	"Content-Encoding",
	"Connection",
	"TE",
	"HTTP2-Settings",
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Connection",
	"Proxy-Authorization",
	"Trailer",
	"Upgrade",
	"Transfer-Encoding",
	"Alternate-Protocol",
	"X-Forwarded-For",
	"Proxy-Connection",
}

func stripHTHHeaders(h *http.Header) {
	for _, header := range HopByHopHeaders {
		h.Del(header)
	}
}

//This is a mock RoundTrip used to intercept responses before they are forwarded by the proxy
func (ri *Interceptor) RoundTrip(req *http.Request) (res *http.Response, err error) {
	//log.Println("Request read by proxy")
	intercept.RLock()
	intercepted := intercept.value
	intercept.RUnlock()
	//log.Println("Preprocessing...")
	backUpURL := req.URL
	req, Id, err := preProcessRequest(req)
	//log.Println("...done")
	if err != nil {
		//TODO
		log.Println(err)
	}
	if intercepted {
		req, err = editRequest(req, Id)
		if err != nil {
			//TODO
			log.Println(err)
		}
		req.URL.Scheme = backUpURL.Scheme
		req.URL.Host = backUpURL.Host
	}
	status.RLock()
	status.ReqResps[Id].parseRequest(req)
	status.RUnlock()

	//Perform the request, but disable compressing.
	//The gzip encoding will be used by the http package
	req.Header.Del("Accept-Encoding")
	//log.Println("Requesting: ", Id)
	res, err = ri.wrappedRT.RoundTrip(req)
	//log.Println("Received response for req: ", Id)
	if err != nil {
		log.Println("Something went wrong trying to contact the server")
		return
	}
	res, err = editResponse(req, res, intercepted, Id)
	status.RLock()
	status.ReqResps[Id].parseResponse(res)
	status.RUnlock()
	return
}
