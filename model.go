package stropt

import (
	"fmt"
	"os"
)

// the helper model for show the help message
type Help struct {
	// this is the helper utility and show the help message
	Help bool `shortcut:"h" name:"help" desc:"show this help message and exit" callback:"Help_"`
}

// the help model for show the version info
type Model struct {
	Help

	// this is the helper utility and show the version info
	Version bool `shortcut:"v" name:"version" desc:"show the version and exit" callback:"Version_"`
}

func init() {
	// regitser all model's callback
	RegisterCallback(CALLBACK_HELP, help)
	RegisterCallback(CALLBACK_VERSION, version)
}

// show the usage on stderr, and exit
func help(stropt *StrOpt, _field Field) (err error) {
	stropt.Usage(os.Stderr)
	os.Exit(1)
	return
}

var (
	// the version info, may override by caller
	ver string
)

// show the version info, may override by caller
func version(stropt *StrOpt, _field Field) (err error) {
	if ver == "" {
		// show the StrOpt version info
		ver = fmt.Sprintf("%v (v%d.%d.%d)", PROJ_NAME, MAJOR, MINOR, MAJOR)
	}

	os.Stdout.WriteString(ver)
	os.Exit(0)
	return
}

func Version(_ver string) {
	ver = _ver
}
