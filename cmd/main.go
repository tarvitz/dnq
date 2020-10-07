package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/jessevdk/go-flags"
)

type versionCommand struct{}

// Main application options
type Main struct {
	Verbose []bool         `short:"v" description:"Show verbose information"`
	Version versionCommand `command:"version" alias:"ver" description:"show command version."`

	// tests only commands
	TestError bool `long:"error" hidden:"true" description:"testing app returns err."`
}

var (
	opts       = &Main{}
	parser     = flags.NewParser(opts, flags.HelpFlag|flags.PassDoubleDash)
	appVersion = "dev"

	// Exit function
	Exit = func (code int) {
		os.Exit(code)
	}
)

func verbosity() int {
	v := 0
	for _, entry := range os.Args[1:] {
		if strings.HasPrefix(entry, "-v") {
			v += strings.Count(entry, "v")
		}
	}
	return v
}

func exit(err error) {
	if err == nil {
		Exit(exitOk)
		// this return makes sense only for testing, due to
		// there's no real system exit from this function, thus far
		// running cmd tests it will continue to follow the code sequence.
		return
	}

	if flagsErr, ok := err.(*flags.Error); ok {
		switch flagsErr.Type {
		case flags.ErrHelp:
			parser.WriteHelp(os.Stdout)
			Exit(exitOk)
		case flags.ErrCommandRequired:
			parser.WriteHelp(os.Stdout)
			Exit(exitCommandRequired)
		default:
			log.Printf("command line parser error: %v", err)
			Exit(exitCodeParserFailure)
		}
	} else {
		Exit(exitCommandError)
	}
}

func (command *versionCommand) Execute(_ []string) (err error) {
	version := fmt.Sprintf(
		"%[1]s version: %[2]s, %[3]s/%[4]s %[5]s",
		appName, appVersion, runtime.GOOS, runtime.GOARCH, runtime.Version())
	fmt.Println(version)
	return
}

func main() {
	_, err := parser.Parse()
	exit(err)
}
