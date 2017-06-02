package apis

type Command struct {
	Channel string
	Action  string
	Args    []string
	Payload []byte
}
