package context_conf

// code-message map
var codes = map[int]string{
	0: "success",
	1: "error",
}

type ErrorCode struct {
	Code int
	Msg  string
}

var (
	SUCCESS = errorCode(0)
	ERROR   = errorCode(1)
)

func errorCode(code int) ErrorCode {
	return ErrorCode{
		Code: code,
		Msg:  codes[code],
	}
}
