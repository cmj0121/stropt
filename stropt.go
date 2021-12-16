package stropt

type StrOpt struct {
}

func New(in interface{}) (stropt *StrOpt, err error) {
	return
}

func MustNew(in interface{}) (stropt *StrOpt) {
	var err error

	if stropt, err = New(in); err != nil {
		panic(err)
	}
	return
}
