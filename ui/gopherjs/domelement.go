// +build js

package main

import (
	"fmt"
	"strings"

	"github.com/empijei/cli/lg"
	"github.com/gopherjs/gopherjs/js"
)

var document *js.Object

type DomElement struct {
	*js.Object
}

func toString(o *js.Object) string {
	if o == nil || o == js.Undefined {
		return ""
	}
	return o.String()
}

func GetElementByID(id string) *DomElement {
	return &DomElement{js.Global.Get(id)}
}

func (de *DomElement) SetTextContent(content string) {
	de.Set("textContent", content)
}
func (de *DomElement) SetTextContentf(format string, args ...interface{}) {
	de.Set("textContent", fmt.Sprintf(format, args...))
}

func (de *DomElement) GetTextContent() string {
	return toString(de.Get("textContent"))
}

func (de *DomElement) ToggleClass(old, new string) {
	oldclasses := strings.Split(toString(de.Get("classList")), " ")
	lg.Debugf("Oldclasses: %v", oldclasses)
	newclasses := make([]string, 0, len(oldclasses)+1)
	var replaced bool
	for _, class := range oldclasses {
		if class == old {
			replaced = true
			if new != "" {
				newclasses = append(newclasses, new)
			}
		} else {
			newclasses = append(newclasses, class)
		}
	}
	if !replaced {
		newclasses = append(newclasses, new)
	}

	lg.Debugf("New classes: %v", newclasses)
	de.Set("classList", strings.Join(newclasses, " "))
}

func (de *DomElement) SetAttribute(name string, value string) {
	de.Call("setAttribute", name, value)
}

func (de *DomElement) SetAttributes(keyvalues map[string]string) {
	for attribute, value := range keyvalues {
		de.SetAttribute(attribute, value)
	}
}

func (de *DomElement) AppendChild(child *DomElement) {
	de.Call("appendChild", child.Object)
}

func createElement(name string) *DomElement {
	return &DomElement{document.Call("createElement", name)}
}
