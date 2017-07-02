package intercept

/*
func a(req *http.Request, in *http.Response) (*http.Response, error) {
	in.Header.Set("X-Wapty-Test", in.Header.Get("X-Wapty-Test")+"FuncA")
	return in, nil
}
func b(req *http.Request, in *http.Response) (*http.Response, error) {
	in.Header.Set("X-Wapty-Test", in.Header.Get("X-Wapty-Test")+"FuncB")
	return in, nil
}
func c(in *http.Request) (*http.Request, error) {
	in.Header.Set("X-Wapty-Test", in.Header.Get("X-Wapty-Test")+"FuncC")
	return in, nil
}
func d(in *http.Request) (*http.Request, error) {
	in.Header.Set("X-Wapty-Test", in.Header.Get("X-Wapty-Test")+"FuncD")
	return in, nil
}

func TestComposeResponseModifier(t *testing.T) {
	var ab = composeResponseModifier(a, b)
	tester := &http.Response{Header: http.Header{}}
	result, _ := ab(nil, tester)
	if actual := result.Header.Get("X-Wapty-Test"); "FuncBFuncA" != actual {
		t.Errorf("composeResponseModifier did not work, expected FuncBFuncA, got " + actual)
	}
	tester.Header.Del("X-Wapty-Test")
	var ba = composeResponseModifier(b, a)
	result, _ = ba(nil, tester)
	if actual := result.Header.Get("X-Wapty-Test"); "FuncAFuncB" != actual {
		t.Errorf("composeResponseModifier did not work, expected FuncAFuncB, got " + actual)
	}
}

func TestPlug(t *testing.T) {
	var tester PlugHandler
	tester.PreProcessRequest(c, true)
	tester.PreProcessRequest(d, true)
	result := &http.Request{Header: http.Header{}}
	result, _ = tester.preModifyRequest(result)
	if actual := result.Header.Get("X-Wapty-Test"); "FuncDFuncC" != actual {
		t.Errorf("PreProcessRequest did not work, expected FuncDFuncC, got " + actual)
	}
}
*/
