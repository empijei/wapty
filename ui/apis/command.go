package apis

type Command struct {
	Channel string
	Action  string
	Args    map[string]string
	Payload []byte
}
