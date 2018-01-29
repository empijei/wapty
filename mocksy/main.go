package mocksy

import "github.com/empijei/cli/lg"

// Main is the main function that starts Mocksy
func Main(_ ...string) {
	const port = ":8082"
	const histDir = "."

	SetHistDir(histDir)

	lg.Infof("Starting mocksy server at %s", port)
	if err := StartServer(port); err != nil {
		panic(err)
	}
}
