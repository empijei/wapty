package mocksy

import "log"

func Main(_ ...string) {
	const port = ":8082"
	const histDir = "."

	SetHistDir(histDir)

	log.Printf("Starting mocksy server at %s\n", port)
	if err := StartServer(port); err != nil {
		panic(err)
	}
}
