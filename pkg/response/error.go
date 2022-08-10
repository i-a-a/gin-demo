package response

// 鸭子类型，自定义error
var (
	_ error = Msg("HelloWorld")
	_ error = Code(67373)
)

type Msg string

func (f Msg) Error() string {
	return string(f)
}

type Code int

func (c Code) Error() string {
	return "Do something"
}
