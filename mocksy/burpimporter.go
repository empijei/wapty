package mocksy

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"log"
)

type Request struct {
	Base64 string `xml:"base64,attr"`
	Value  []byte `xml:",chardata"`
}

func (r Request) Bytes() []byte {
	return b64Able(r).Bytes()
}

type Response struct {
	Base64 string `xml:"base64,attr"`
	Value  []byte `xml:",chardata"`
}

func (r Response) Bytes() []byte {
	return b64Able(r).Bytes()
}

type b64Able struct {
	Base64 string
	Value  []byte
}

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

type Host struct {
	Ip    string `xml:"ip,attr"`
	Value string
}

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

type Items struct {
	Items []Item `xml:"item"`
}

func BurpImport(r io.Reader) (*Items, error) {
	dec := xml.NewDecoder(r)
	var itm Items
	err := dec.Decode(&itm)
	if err != nil {
		//wrapping errors is good practice
		return nil, fmt.Errorf("mocksy: cannot import status: %s", err.Error())
	}
	return &itm, nil
}
