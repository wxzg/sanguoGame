package common

type Result struct {
	Code int `json:"code"`
	Errmsg string `json:"errmsg"`
	Data interface{} `json:"data"`
}

func Error(code int, msg string)  (r *Result) {
	r = &Result{}
	r.Code = code
	r.Errmsg = msg
	return
}

func Success(code int, data interface{}) (r *Result)  {
	r = &Result{}
	r.Code = code
	r.Data = data
	return
}