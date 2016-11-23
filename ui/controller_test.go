package ui

type MockClient struct {
}

func NewMockClient() {

}
func (m *MockClient) Receive() Command {
	return Command{}
}
