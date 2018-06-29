package message

import "github.com/zhangmingfeng/minres/controllers/base/message"

const (
	FileNameIsEmpty    = 10001
	FileSizeIsInvalid  = 10002
	FileTimeIsInvalid  = 10003
	ChunkSizeIsInvalid = 10004
)

type ParamsRequest struct {
	FileName  string `json:"fileName"`
	FileSize  int    `json:"fileSize"`
	FileTime  int    `json:"fileTime"`
	ChunkSize int    `json:"chunkSize"`
}

type ParamsResponse struct {
	message.BaseResponse
	Token     string `json:"token"`
	UploadUrl string `json:"uploadUrl"`
	Params    Params `json:"params"`
}

type Params struct {
	Loaded     int  `json:"loaded"`
	IsFinished bool `json:"isFinished"`
}

func NewParamsRequest() *ParamsRequest {
	return &ParamsRequest{}
}

func NewParamsResponse() *ParamsResponse {
	return &ParamsResponse{
		BaseResponse: message.BaseResponse{
			Code: 200,
		},
		Params: Params{},
	}
}
