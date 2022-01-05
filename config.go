package stropt

const (
	// the project name
	PROJ_NAME = "stropt"
	// the version meta
	MAJOR = 0 // bump when API changed
	MINOR = 2 // bump when new feature add
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
	// the pre-defined choice of the input
	KEY_CHOICE = "choice"
	// the default value
	KEY_DEFAULT = "default"

	// the attribute of field
	KEY_ATTR      = "attr"
	KEY_ATTR_FLAG = "flag"
)

// pre-defined tag used in stropt
var (
	// the field should be ignored in stropt
	TAG_IGNORE = "-"
)
