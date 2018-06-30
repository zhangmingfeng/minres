package message

import "github.com/zhangmingfeng/minres/controllers/base/message"

const (
	FileNameIsEmpty    = 10001
	FileSizeIsInvalid  = 10002
	FileTimeIsInvalid  = 10003
	ChunkSizeIsInvalid = 10004
	TokenIsEmpty       = 10005
	TokenIsInvalid     = 10006
)

type ParamsRequest struct {
	FileName  string `json:"fileName,omitempty"`
	FileSize  int    `json:"fileSize,omitempty"`
	FileTime  int    `json:"fileTime,omitempty"`
	ChunkSize int    `json:"chunkSize,omitempty"`
}

type ParamsResponse struct {
	message.BaseResponse
	Token     string `json:"token,omitempty"`
	UploadUrl string `json:"uploadUrl,omitempty"`
	Params    Params `json:"params,omitempty"`
}

type Params struct {
	Loaded     int    `json:"loaded,omitempty"`
	Key        string `json:"key,omitempty"`
	IsFinished bool   `json:"isFinished,omitempty"`
}

type UploadRequest struct {
	Token  string `json:"token,omitempty"`
	Chunk  int    `json:"chunk,omitempty"`
	Chunks int    `json:"chunks,omitempty"`
}

type UploadResponse struct {
	message.BaseResponse
	IsFinished string `json:"isFinished,omitempty"`
	Loaded     int    `json:"loaded,omitempty"`
	Chunk      int    `json:"chunk,omitempty"`
	Chunks     int    `json:"chunks,omitempty"`
	File       File   `json:"file,omitempty"`
}

type File struct {
	Key string `json:"key,omitempty"`
	Url string `json:"url,omitempty"`
}

func NewParamsRequest() *ParamsRequest {
	return &ParamsRequest{}
}

func NewUploadRequest() *UploadRequest {
	return &UploadRequest{}
}

func NewParamsResponse() *ParamsResponse {
	return &ParamsResponse{
		BaseResponse: message.BaseResponse{
			Code: 200,
			Msg:  "success",
		},
	}
}
func NewUploadResponse() *UploadResponse {
	return &UploadResponse{
		BaseResponse: message.BaseResponse{
			Code: 200,
			Msg:  "success",
		},
	}
}
