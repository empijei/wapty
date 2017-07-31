package mocksy

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"log"
)

// Response is struct used to deserialize burp XML. It contains text data and an attribute telling
// if the content is base64 or not.
type Request struct {
	Base64 string `xml:"base64,attr"`
	Value  []byte `xml:",chardata"`
}

// Bytes returns the []byte read by the XML decoder
func (r Request) Bytes() []byte {
	return b64Able(r).Bytes()
}

// Response is struct used to deserialize burp XML. It contains text data and an attribute telling
// if the content is base64 or not.
type Response struct {
	Base64 string `xml:"base64,attr"`
	Value  []byte `xml:",chardata"`
}

// Bytes returns the []byte read by the XML decoder
func (r Response) Bytes() []byte {
	return b64Able(r).Bytes()
}

// this is just a sort of interface for fields used to not repeat code for the
// Bytes function
type b64Able struct {
	Base64 string
	Value  []byte
}

// Bytes returns the []byte read by the XML decoder
func (r b64Able) Bytes() []byte {
	if r.Base64 == "true" {
		value, err := base64.StdEncoding.DecodeString(string(r.Value))
		if err != nil {
			//TODO handle more gently
			log.Fatal(err)
		}
		return value
	}
	return r.Value
}

// Host is used to deserialize burp saved XML. It contains text data and an attribute "ip".
type Host struct {
	Ip    string `xml:"ip,attr"`
	Value string `xml:",chardata"`
}

// Item is used to deserialize burp saved XML. It is the element containing all the information
// about a single request.
type Item struct {
	Time           string   `xml:"time"`
	Url            string   `xml:"url"`
	Request        Request  `xml:"request"`
	Host           Host     `xml:"host"`
	Port           string   `xml:"port"`
	Protocol       string   `xml:"protocol"`
	Method         string   `xml:"method"`
	Path           string   `xml:"path"`
	Extension      string   `xml:"extension"`
	Status         string   `xml:"status"`
	ResponseLength string   `xml:"responselength"`
	Response       Response `xml:"response"`
	Mimetype       string   `xml:"mimetype"`
	Comment        string   `xml:"comment"`
}

// Items is used to deserialize burp saved XML. It is the root element, containing a list of Item's.
type Items struct {
	Items []Item `xml:"item"`
}

// BurpImport reads a "saved requests" file from r both in base64 and cleartext form
func BurpImport(r io.Reader) (*Items, error) {
	dec := xml.NewDecoder(r)
	var itm Items
	err := dec.Decode(&itm)
	if err != nil {
		// wrapping errors is good practice
		return nil, fmt.Errorf("mocksy: cannot import status: %s", err.Error())
	}
	return &itm, nil
}
