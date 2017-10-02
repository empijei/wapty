package intercept

import (
	"net/http"

	"github.com/empijei/wapty/cli/lg"
)

//Remove trailers?
//https://github.com/squid-cache/squid/blob/master/src/http/RegisteredHeadersHash.cci

// HopByHopHeaders is a list of all the HTTP headers that are stripped away by
// the proxy
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

// Interceptor is a struct that respects the net.RoundTripper interface and just wraps
// the original http.RoundTripper
type Interceptor struct {
	wrappedRT http.RoundTripper
}

// RoundTrip is a mock RoundTrip used to intercept requests and responses
// before they are forwarded by the proxy.
func (ri *Interceptor) RoundTrip(req *http.Request) (res *http.Response, err error) {
	//This first part is dedicated to the REQUESTS
	intercepted := intercept.value()
	backUpURL := req.URL
	req, ID, err := preProcessRequest(req)
	if err != nil {
		//TODO handle possible autodrop
		//TODO other errors
		lg.Error(err)
	}
	if intercepted {
		var editedReq *http.Request
		editedReq, res, err = editRequest(req, ID)
		if err != nil {
			//TODO
			lg.Error(err)
		}
		if editedReq != nil {
			req = editedReq
			req.URL.Scheme = backUpURL.Scheme
			req.URL.Host = backUpURL.Host
		}
	}

	status.RLock()
	status.ReqResps[ID].parseRequest(req)
	status.RUnlock()
	if res != nil {
		//TODO Adding dropped responses could be avoided.
		status.addResponse(ID, res)
		return
	}

	//This second part works on the RESPONSES
	//Perform the request, but disable compressing.
	//The gzip encoding should be used by the http package transparently
	req.Header.Del("Accept-Encoding")
	res, err = ri.wrappedRT.RoundTrip(req)
	if err != nil {
		lg.Error("Something went wrong trying to contact the server")
		//TODO return a fake response containing the error message
		res = GenerateResponse("Error", "Error in performing the request: "+err.Error(), 500)
		return
	}
	res = preProcessResponse(req, res, ID)
	if intercepted {
		res, err = editResponse(req, res, ID)
	}
	status.RLock()
	status.ReqResps[ID].parseResponse(res)
	status.RUnlock()
	return
}
