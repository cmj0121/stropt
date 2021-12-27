package stropt

import (
	"os"
)

// the helper model for show the help message
type Help struct {
	// this is the helper utility and show the help message
	Help bool `shortcut:"h" name:"help" desc:"show this help message" callback:"_help"`
}

func init() {
	// regitser all model's callback
	RegisterCallback(CALLBACK_HELP, help)
}

// show the usage on stderr, and exit
func help(stropt *StrOpt, _field Field) {
	stropt.Usage(os.Stderr)
	os.Exit(1)
}
