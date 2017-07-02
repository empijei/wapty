package intercept

/*
// Plugin is the instance of the resulting composition of all the plugins. This is likely
// to be changed in future versions so please do not use/modify this.
var Plugin PlugHandler

// PlugHandler is the resulting composition of all the plugins. This is likely
// to be changed in future versions so please do not use/modify this.
type PlugHandler struct {
	sync.Mutex
	used bool
	//Called if request is intercepted and before the buffer is sento to the UI for editing
	preModifyRequest RequestModifier
	//called on every request before preModifyRequest
	alwaysModifyRequest RequestModifier
	//Called if request is intercepted and after the buffer has been modified
	postModifyRequest RequestModifier
	//Called if response is intercepted and before the buffer is sento to the UI for editing
	preModifyResponse ResponseModifier
	//called on every response before preModifyRequest
	alwaysModifyRespone ResponseModifier
	//Called if response is intercepted and after the buffer has been modified
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
*/
