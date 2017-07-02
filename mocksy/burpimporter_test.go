package mocksy

import (
	"bytes"
	"fmt"
	"testing"
)

const data = `
<items>
<item>
    <time>Wed May 31 20:25:16 CEST 2017</time>
    <url><![CDATA[http://www.ansa.it/sito/img/ico_spread_dwn.png]]></url>
    <host ip="194.244.5.206">www.ansa.it</host>
    <port>80</port>
    <protocol>http</protocol>
    <method><![CDATA[GET]]></method>
    <path><![CDATA[/sito/img/ico_spread_dwn.png]]></path>
    <extension>png</extension>
    <request base64="false"><![CDATA[GET /sito/img/ico_spread_dwn.png HTTP/1.1
Host: www.ansa.it
User-Agent: Mozilla/5.0 (X11; Linux x86_64; rv:53.0) Gecko/20100101 Firefox/53.0
Accept: */*
Accept-Language: en-US,en;q=0.5
Referer: http://www.ansa.it/
Connection: close

]]></request>
    <status>302</status>
    <responselength>395</responselength>
    <mimetype>HTML</mimetype>
    <response base64="false"><![CDATA[HTTP/1.1 302 Moved Temporarily
Content-Type: text/html
Location: http://wwwrb.ansa.it/sito/img/ico_spread_dwn.png
X-Node: www3
Vary: Accept-Encoding
Date: Wed, 31 May 2017 18:25:16 GMT
Connection: close
age: 0
Content-Length: 154

<html>
<head><title>302 Found</title></head>
<body bgcolor="white">
<center><h1>302 Found</h1></center>
<hr><center>nginx</center>
</body>
</html>
]]></response>
    <comment></comment>
 </item>
  <item>
    <time>Tue Jun 06 22:34:15 CEST 2017</time>
    <url><![CDATA[https://www.google.it/gen_204?atyp=i&ct=&cad=&vet=10ahUKEwiesrjXiKrUAhUD2CwKHSBNBOYQsmQIDw..s&ei=xRE3Wd67J4OwswGgmpGwDg&zx=1496781254993]]></url>
    <host ip="216.58.198.3">www.google.it</host>
    <port>443</port>
    <protocol>https</protocol>
    <method><![CDATA[GET]]></method>
    <path><![CDATA[/gen_204?atyp=i&ct=&cad=&vet=10ahUKEwiesrjXiKrUAhUD2CwKHSBNBOYQsmQIDw..s&ei=xRE3Wd67J4OwswGgmpGwDg&zx=1496781254993]]></path>
    <extension>null</extension>
    <request base64="true"><![CDATA[R0VUIC9nZW5fMjA0P2F0eXA9aSZjdD0mY2FkPSZ2ZXQ9MTBhaFVLRXdpZXNyalhpS3JVQWhVRDJDd0tIU0JOQk9ZUXNtUUlEdy4ucyZlaT14UkUzV2Q2N0o0T3dzd0dnbXBHd0RnJnp4PTE0OTY3ODEyNTQ5OTMgSFRUUC8xLjENCkhvc3Q6IHd3dy5nb29nbGUuaXQNClVzZXItQWdlbnQ6IE1vemlsbGEvNS4wIChYMTE7IExpbnV4IHg4Nl82NDsgcnY6NTMuMCkgR2Vja28vMjAxMDAxMDEgRmlyZWZveC81My4wDQpBY2NlcHQ6ICovKg0KQWNjZXB0LUxhbmd1YWdlOiBlbi1VUyxlbjtxPTAuNQ0KUmVmZXJlcjogaHR0cHM6Ly93d3cuZ29vZ2xlLml0Lw0KQ29va2llOiBOSUQ9MTA1PVVNNXQzZlpsMTBYbWctaWgxS0R4VjNzWFEyaTlvdU5kYllMMTk0MkxvVUJmVXJFaGYxVDBodndBQ1Bhb3dsbGY0Nk8zeTdSc3d2V0s2cXEta0dmNFZTcWkzRDQ2SW5BNmV1MC1uVDNadS00eGpBalNxOEE4VDExVzJUeTRqdmFhOyBDT05TRU5UPVdQLjI2MTBhYjsgR1o9Wj0xDQpETlQ6IDENCkNvbm5lY3Rpb246IGNsb3NlDQoNCg==]]></request>
    <status>204</status>
    <responselength>268</responselength>
    <mimetype></mimetype>
    <response base64="true"><![CDATA[SFRUUC8xLjEgMjA0IE5vIENvbnRlbnQNCkNvbnRlbnQtVHlwZTogdGV4dC9odG1sOyBjaGFyc2V0PVVURi04DQpEYXRlOiBUdWUsIDA2IEp1biAyMDE3IDIwOjM0OjE1IEdNVA0KU2VydmVyOiBnd3MNCkNvbnRlbnQtTGVuZ3RoOiAwDQpYLVhTUy1Qcm90ZWN0aW9uOiAxOyBtb2RlPWJsb2NrDQpYLUZyYW1lLU9wdGlvbnM6IFNBTUVPUklHSU4NCkFsdC1TdmM6IHF1aWM9Ijo0NDMiOyBtYT0yNTkyMDAwOyB2PSIzOCwzNywzNiwzNSINCkNvbm5lY3Rpb246IGNsb3NlDQoNCg==]]></response>
    <comment></comment>
  </item>
  </items>
 `

func TestBurpImporter(t *testing.T) {
	testbuf := bytes.NewBuffer([]byte(data))
	itm, err := BurpImport(testbuf)
	if err != nil {
		t.Fatal(err)
	}
	if len(itm.Items) != 2 {
		t.Fatal(fmt.Errorf("Expected items to have length 2, but have length %d", len(itm.Items)))
	}
	// ======================= Item 0
	item0 := itm.Items[0]
	assertEqual(t, item0.Time, "Wed May 31 20:25:16 CEST 2017")
	assertEqual(t, item0.Url, "http://www.ansa.it/sito/img/ico_spread_dwn.png")
	assertEqual(t, item0.Host.Value, "www.ansa.it")
	assertEqual(t, item0.Host.Ip, "194.244.5.206")
	assertEqual(t, item0.Port, "80")
	assertEqual(t, item0.Protocol, "http")
	assertEqual(t, item0.Method, "GET")
	assertEqual(t, item0.Path, "/sito/img/ico_spread_dwn.png")
	assertEqual(t, item0.Extension, "png")
	assertEqual(t, item0.Request.Base64, "false")
	assertEqualSlice(t, item0.Request.Value, []byte(`GET /sito/img/ico_spread_dwn.png HTTP/1.1
Host: www.ansa.it
User-Agent: Mozilla/5.0 (X11; Linux x86_64; rv:53.0) Gecko/20100101 Firefox/53.0
Accept: */*
Accept-Language: en-US,en;q=0.5
Referer: http://www.ansa.it/
Connection: close

`))
	assertEqual(t, item0.Status, "302")
	assertEqual(t, item0.ResponseLength, "395")
	assertEqual(t, item0.Mimetype, "HTML")
	assertEqual(t, item0.Response.Base64, "false")
	assertEqualSlice(t, item0.Response.Value, []byte(`HTTP/1.1 302 Moved Temporarily
Content-Type: text/html
Location: http://wwwrb.ansa.it/sito/img/ico_spread_dwn.png
X-Node: www3
Vary: Accept-Encoding
Date: Wed, 31 May 2017 18:25:16 GMT
Connection: close
age: 0
Content-Length: 154

<html>
<head><title>302 Found</title></head>
<body bgcolor="white">
<center><h1>302 Found</h1></center>
<hr><center>nginx</center>
</body>
</html>
`))
	assertEqual(t, item0.Comment, "")

	// ======================= Item 1
	item1 := itm.Items[1]
	assertEqual(t, item1.Time, "Tue Jun 06 22:34:15 CEST 2017")
	assertEqual(t, item1.Url, "https://www.google.it/gen_204?atyp=i&ct=&cad=&vet=10ahUKEwiesrjXiKrUAhUD2CwKHSBNBOYQsmQIDw..s&ei=xRE3Wd67J4OwswGgmpGwDg&zx=1496781254993")
	assertEqual(t, item1.Host.Value, "www.google.it")
	assertEqual(t, item1.Host.Ip, "216.58.198.3")
	assertEqual(t, item1.Port, "443")
	assertEqual(t, item1.Protocol, "https")
	assertEqual(t, item1.Method, "GET")
	assertEqual(t, item1.Path, "/gen_204?atyp=i&ct=&cad=&vet=10ahUKEwiesrjXiKrUAhUD2CwKHSBNBOYQsmQIDw..s&ei=xRE3Wd67J4OwswGgmpGwDg&zx=1496781254993")
	assertEqual(t, item1.Extension, "null")
	assertEqual(t, item1.Request.Base64, "true")
	assertEqualSlice(t, item1.Request.Value, []byte(`R0VUIC9nZW5fMjA0P2F0eXA9aSZjdD0mY2FkPSZ2ZXQ9MTBhaFVLRXdpZXNyalhpS3JVQWhVRDJDd0tIU0JOQk9ZUXNtUUlEdy4ucyZlaT14UkUzV2Q2N0o0T3dzd0dnbXBHd0RnJnp4PTE0OTY3ODEyNTQ5OTMgSFRUUC8xLjENCkhvc3Q6IHd3dy5nb29nbGUuaXQNClVzZXItQWdlbnQ6IE1vemlsbGEvNS4wIChYMTE7IExpbnV4IHg4Nl82NDsgcnY6NTMuMCkgR2Vja28vMjAxMDAxMDEgRmlyZWZveC81My4wDQpBY2NlcHQ6ICovKg0KQWNjZXB0LUxhbmd1YWdlOiBlbi1VUyxlbjtxPTAuNQ0KUmVmZXJlcjogaHR0cHM6Ly93d3cuZ29vZ2xlLml0Lw0KQ29va2llOiBOSUQ9MTA1PVVNNXQzZlpsMTBYbWctaWgxS0R4VjNzWFEyaTlvdU5kYllMMTk0MkxvVUJmVXJFaGYxVDBodndBQ1Bhb3dsbGY0Nk8zeTdSc3d2V0s2cXEta0dmNFZTcWkzRDQ2SW5BNmV1MC1uVDNadS00eGpBalNxOEE4VDExVzJUeTRqdmFhOyBDT05TRU5UPVdQLjI2MTBhYjsgR1o9Wj0xDQpETlQ6IDENCkNvbm5lY3Rpb246IGNsb3NlDQoNCg==`))
	assertEqual(t, item1.Status, "204")
	assertEqual(t, item1.ResponseLength, "268")
	assertEqual(t, item1.Mimetype, "")
	assertEqual(t, item1.Response.Base64, "true")
	assertEqualSlice(t, item1.Response.Value, []byte(`SFRUUC8xLjEgMjA0IE5vIENvbnRlbnQNCkNvbnRlbnQtVHlwZTogdGV4dC9odG1sOyBjaGFyc2V0PVVURi04DQpEYXRlOiBUdWUsIDA2IEp1biAyMDE3IDIwOjM0OjE1IEdNVA0KU2VydmVyOiBnd3MNCkNvbnRlbnQtTGVuZ3RoOiAwDQpYLVhTUy1Qcm90ZWN0aW9uOiAxOyBtb2RlPWJsb2NrDQpYLUZyYW1lLU9wdGlvbnM6IFNBTUVPUklHSU4NCkFsdC1TdmM6IHF1aWM9Ijo0NDMiOyBtYT0yNTkyMDAwOyB2PSIzOCwzNywzNiwzNSINCkNvbm5lY3Rpb246IGNsb3NlDQoNCg==`))
	assertEqual(t, item1.Comment, "")
}

func assertEqual(t *testing.T, a, b interface{}) {
	if a != b {
		t.Fatal("Expected:\n\t^", b, "$\nbut got:\n\t^", a, "$")
	}
}

func assertEqualSlice(t *testing.T, a, b []byte) {
	if !bytes.Equal(a, b) {
		t.Fatal("Expected:\n\t^", string(b), "$\nbut got:\n\t^", string(a), "$")
	}
}
