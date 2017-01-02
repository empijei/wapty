package intercept

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"net/http"
)

const templateResponse = `
<html>
<head>
<title>{{.Title}}</title>
</head>
<body>
<!-- TODO WAPTY LOGO -->
<h1>{{.Content}}</h1>
</body>
</html>
`

func GenerateResponse(title, content string, status int) *http.Response {
	t, _ := template.New("Generated Response").Parse(templateResponse)
	data := struct {
		Title, Content string
	}{
		title,
		content,
	}
	body := bytes.NewBuffer(nil)
	_ = t.Execute(body, data)
	res := &http.Response{}
	res.ContentLength = int64(body.Len())
	res.Body = ioutil.NopCloser(body)
	res.StatusCode = status
	res.Header = http.Header{}
	res.Header.Set("X-WAPTY-Status", title)
	return res
}
