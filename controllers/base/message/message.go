package message

type BaseResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}
