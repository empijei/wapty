package cli

import (
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/empijei/wapty/cli/lg"
)

func init() {

	// Setup fallback version and commit in case wapty wasn't "properly" compiled
	if len(Version) == 0 {
		Version = "Unknown, please compile wapty with 'make'"
	}
	if len(Commit) == 0 {
		Commit = "Unknown, please compile wapty with 'make'"
	}
	AddCommand(cmdVersion)
}

const banner = `                     _         
__      ____ _ _ __ | |_ _   _ 
\ \ /\ / / _' | '_ \| __| | | |   Version: {{.Version}}
 \ V  V / (_| | |_) | |_| |_| |   Commit:  {{.Commit}}
  \_/\_/ \__,_| .__/ \__|\__, |   Build:   {{.Build}}
              |_|        |___/    
 
`

var (
	//Version is taken by the build flags, represent current version as
	//<major>.<minor>.<patch>
	Version string

	//Commit is the output of `git rev-parse HEAD` at the moment of the build
	Commit string

	//Build contains info about the scope of the build and should either be Debug or Release
	Build string = "Debug"
)

var cmdVersion = &Cmd{
	Name: "version",
	Run: func(_ ...string) {
		fmt.Printf("Version: %s\nCommit: %s\n", Version, Commit)
	},
	UsageLine: "version",
	Short:     "print version and exit",
	Long:      "print version and exit",
}

func Printbanner() {
	tmpl := template.New("banner")
	template.Must(tmpl.Parse(banner))
	_ = tmpl.Execute(os.Stderr, struct{ Version, Commit, Build string }{Version, Commit, Build})
}

func Init() {
	if Build == "Release" {
		lg.CurLevel = lg.Level_Info
		lg.SetFlags(log.Ltime)
	}

	stderrinfo, err := os.Stderr.Stat()
	if err == nil && stderrinfo.Mode()&os.ModeCharDevice == 0 {
		// Output is a pipe, turn off colors
		lg.Color = false
	} else {
		// Output is to terminal, print banner
		Printbanner()
	}

	if len(os.Args) > 1 {
		//read the first argument
		directive := os.Args[1]
		if len(os.Args) > 2 {
			//shift parameters left, but keep argv[0]
			os.Args = append(os.Args[:1], os.Args[2:]...)
		} else {
			os.Args = os.Args[:1]
		}
		command, err := FindCommand(directive)
		if err == nil {
			callCommand(command)
		} else {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			fmt.Fprintln(os.Stderr, "Available commands are:\n")
			for _, cmd := range WaptyCommands {
				fmt.Fprintln(os.Stderr, "\t"+cmd.Name+"\n\t\t"+cmd.Short)
			}
			fmt.Fprintln(os.Stderr, "\nDefault command is: ", DefaultCommand.Name)
		}
	} else {
		callCommand(DefaultCommand)
	}
}
