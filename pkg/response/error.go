package response

var (
	_ error = Action{}
	_ error = String("")
)

type Action struct {
	Code int
	Msg  string
}

func (a Action) Error() string {
	return a.Msg
}

type String string

func (f String) Error() string {
	return string(f)
}
