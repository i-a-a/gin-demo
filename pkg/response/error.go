package response

// 鸭子类型，自定义error
var (
	_ error = Msg("HelloWorld")
	_ error = Code(0)
)

type Msg string

func (f Msg) Error() string {
	return string(f)
}

type Code int

func (c Code) Error() string {
	s, _ := codes[int(c)]
	return s
}

var codes = map[int]string{}

func ExtendCodeMap(m map[int]string) {
	for k, v := range m {
		codes[k] = v
	}
}
