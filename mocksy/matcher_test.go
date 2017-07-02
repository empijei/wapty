package mocksy

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

const matcherTestData = `
<items burpVersion="1.7.23" exportTime="Tue Jun 06 22:35:27 CEST 2017">
  <item>
    <time>Tue Jun 06 22:34:15 CEST 2017</time>
    <url><![CDATA[https://localhost/]]></url>
    <host ip="0.0.0.0">localhost</host>
    <port>443</port>
    <protocol>http</protocol>
    <method><![CDATA[POST]]></method>
    <path><![CDATA[/]]></path>
    <extension>null</extension>
    <request base64="false"><![CDATA[PING]]></request>
    <status>200</status>
    <responselength>4</responselength>
    <mimetype></mimetype>
    <response base64="false"><![CDATA[PONG]]></response>
    <comment></comment>
  </item>
  <item>
    <time>Tue Jun 06 22:34:15 CEST 2017</time>
    <url><![CDATA[https://localhost/]]></url>
    <host ip="0.0.0.0">localhost</host>
    <port>443</port>
    <protocol>http</protocol>
    <method><![CDATA[GET]]></method>
    <path><![CDATA[/]]></path>
    <extension>null</extension>
    <request base64="false"></request>
    <status>200</status>
    <responselength>4</responselength>
    <mimetype></mimetype>
    <response base64="false"><![CDATA[GET]]></response>
    <comment></comment>
  </item>
  <item>
    <time>Tue Jun 06 22:34:15 CEST 2017</time>
    <url><![CDATA[https://localhost:8082/asd]]></url>
    <host ip="0.0.0.0">localhost</host>
    <port>8082</port>
    <protocol>http</protocol>
    <method><![CDATA[POST]]></method>
    <path><![CDATA[/asd]]></path>
    <extension>null</extension>
    <request base64="false"><![CDATA[PING]]></request>
    <status>200</status>
    <responselength>4</responselength>
    <mimetype></mimetype>
    <response base64="false"><![CDATA[8082]]></response>
    <comment></comment>
  </item>
  <item>
    <time>Tue Jun 06 22:34:15 CEST 2017</time>
    <url><![CDATA[https://localhost:8082/asd]]></url>
    <host ip="0.0.0.0">localhost</host>
    <port>8083</port>
    <protocol>http</protocol>
    <method><![CDATA[POST]]></method>
    <path><![CDATA[/asd]]></path>
    <extension>null</extension>
    <request base64="false"><![CDATA[PINGA]]></request>
    <status>200</status>
    <responselength>5</responselength>
    <mimetype></mimetype>
    <response base64="false"><![CDATA[PONGA]]></response>
    <comment></comment>
  </item>
</items>`

func init() {
	err := LoadResponsesFrom(strings.NewReader(matcherTestData))
	if err != nil {
		panic(err)
	}
}

// httpBody is a trivial ReadCloser to pass as Body to http.Request
type httpBody struct {
	io.Reader
}

func (h httpBody) Close() error { return nil }

func TestMatcher(t *testing.T) {
	u, _ := url.Parse("http://localhost/")
	req := http.Request{
		Method:        "GET",
		URL:           u,
		Proto:         "HTTP/1.0",
		Body:          httpBody{strings.NewReader(`PING`)},
		ContentLength: 4,
		Host:          "localhost",
	}
	resp := FindMatching(&req)

	if !bytes.Equal(resp.Value, []byte("PONG")) {
		t.Fatal("Expected response 'PONG' but got", string(resp.Value))
	}
}
