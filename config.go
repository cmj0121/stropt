package stropt

const (
	// the project name
	PROJ_NAME = "stropt"
	// the version meta
	MAJOR = 0 // bump when API changed
	MINOR = 1 // bump when new feature add
	MACRO = 0 // bump when bug fixed
)

// pre-defined key of the tag
var (
	// the shortcut of the field
	KEY_SHORTCUT = "shortcut"
	// the name of the field
	KEY_NAME = "name"
	// the description of the field
	KEY_DESC = "desc"
	// the callback function, may local or global
	KEY_CALLBACK = "callback"
	// the default value
	KEY_DEFAULT = "default"
)

// pre-defined tag used in stropt
var (
	// the field should be ignored in stropt
	TAG_IGNORE = "-"
)
