package mocksy

import "log"

func Main() {
	const port = ":8082"

	log.Printf("Starting mocksy server at %s\n", port)
	if err := StartServer(port); err != nil {
		panic(err)
	}
}
