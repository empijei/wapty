package ui

import "github.com/empijei/wapty/ui/apis"

type MockClient struct {
}

func NewMockClient() {

}
func (m *MockClient) Receive() apis.Command {
	return apis.Command{}
}
