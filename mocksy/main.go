package mocksy

import "github.com/empijei/wapty/cli/lg"

func Main(_ ...string) {
	const port = ":8082"
	const histDir = "."

	SetHistDir(histDir)

	lg.Infof("Starting mocksy server at %s\n", port)
	if err := StartServer(port); err != nil {
		panic(err)
	}
}
