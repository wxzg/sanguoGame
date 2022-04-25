package common

type MyError struct {
	code int
	msg string
}

func (m *MyError) Error() string {
	return m.msg
}

func (m *MyError) Code() int {
	return m.code
}

// 定义error
func New(code int,msg string) error  {
	return &MyError{
		code: code,
		msg: msg,
	}
}