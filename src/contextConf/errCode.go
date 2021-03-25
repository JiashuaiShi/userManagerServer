package contextConf

type ErrorCode struct {
	Code int
	Msg  string
}

func newcode(code int) *ErrorCode {
	return &ErrorCode{
		Code: code,
		Msg:  codes[code],
	}
}

var codes = map[int]string{
	0: "success",
	1: "error",
}

var (
	SUCCESS = newcode(0)
	ERROR   = newcode(1)
)
