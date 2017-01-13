package intercept

import (
	"net/http"
	"sync"
)

var Plugin PlugHandler

type PlugHandler struct {
	sync.Mutex
	used               bool
	preModifyRequest   RequestModifier
	postModifyRequest  RequestModifier
	preModifyResponse  ResponseModifier
	postModifyResponse ResponseModifier
}
type RequestModifier func(*http.Request) (*http.Request, error)

func (p *PlugHandler) PreProcessRequest(rm RequestModifier, pre bool) {
	p.Lock()
	defer p.Unlock()
	addRequestModifier(&p.preModifyRequest, rm, pre)
}

func (p *PlugHandler) PostProcessRequest(rm RequestModifier, pre bool) {
	p.Lock()
	defer p.Unlock()
	addRequestModifier(&p.postModifyRequest, rm, pre)
}

func addRequestModifier(field *RequestModifier, rm RequestModifier, pre bool) {
	if *field == nil {
		*field = rm
	} else {
		if pre {
			*field = composeRequestModifier(*field, rm)
		} else {
			*field = composeRequestModifier(rm, *field)
		}
	}
}

func composeRequestModifier(a RequestModifier, b RequestModifier) RequestModifier {
	return func(r *http.Request) (*http.Request, error) {
		out, err := b(r)
		if err != nil {
			return nil, err
		}
		return a(out)
	}
}

type ResponseModifier func(*http.Request, *http.Response) (*http.Response, error)

func (p *PlugHandler) PreProcessResponse(rm ResponseModifier, pre bool) {
	p.Lock()
	defer p.Unlock()
	addResponseModifier(&p.preModifyResponse, rm, pre)
}

func (p *PlugHandler) PostProcessResponse(rm ResponseModifier, pre bool) {
	p.Lock()
	defer p.Unlock()
	addResponseModifier(&p.postModifyResponse, rm, pre)
}

func addResponseModifier(field *ResponseModifier, rm ResponseModifier, pre bool) {
	if *field == nil {
		*field = rm
	} else {
		if pre {
			*field = composeResponseModifier(*field, rm)
		} else {
			*field = composeResponseModifier(rm, *field)
		}
	}
}

func composeResponseModifier(a ResponseModifier, b ResponseModifier) ResponseModifier {
	return func(req *http.Request, in *http.Response) (*http.Response, error) {
		out, err := b(req, in)
		if err != nil {
			return nil, err
		}
		return a(req, out)
	}
}
