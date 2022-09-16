package carrot

// 鸭子类型，这是一个error
var (
	_ error = Msg("")
)

type Msg string

func (msg Msg) Error() string {
	return string(msg)
}
