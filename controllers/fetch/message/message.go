package message

type Request struct {
	Fid      string `json:"fid,omitempty"`
	Width    int    `json:"w,omitempty"`
	Height   int    `json:"h,omitempty"`
	Mode     string `json:"m,omitempty"`
	Download bool   `json:"dl,omitempty"`
}

func NewRequest() *Request {
	return &Request{}
}
