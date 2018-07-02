package message

import (
	"github.com/zhangmingfeng/minres/controllers/base/message"
	"github.com/zhangmingfeng/minres/plugins/seaweedfs"
)

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
	FileGroup string `json:"fileGroup,omitempty"`
	FileSize  int64  `json:"fileSize,omitempty"`
	FileTime  int64  `json:"fileTime,omitempty"`
	ChunkSize int64  `json:"chunkSize,omitempty"`
}

type ParamsResponse struct {
	message.BaseResponse
	Token     string `json:"token,omitempty"`
	UploadUrl string `json:"uploadUrl,omitempty"`
	Loaded    int64  `json:"loaded,omitempty"`
	ChunkSize int64  `json:"chunkSize,omitempty"`
	Chunk     int    `json:"chunk,omitempty"`
	Chunks    int    `json:"chunks,omitempty"`
	FileGroup string `json:"fileGroup,omitempty"`
}

type UploadRequest struct {
	Token string `json:"token,omitempty"`
	File  string `json:"file,omitempty"`
	Chunk int    `json:"chunk,omitempty"`
}

type UploadResponse struct {
	message.BaseResponse
	IsFinished bool  `json:"isFinished,omitempty"`
	Loaded     int64 `json:"loaded,omitempty"`
	Chunk      int   `json:"chunk,omitempty"`
	Chunks     int   `json:"chunks,omitempty"`
	File       File  `json:"file,omitempty"`
}

type File struct {
	Fid string `json:"fid,omitempty"`
	Url string `json:"url,omitempty"`
}

type TokenData struct {
	Fid        string `json:"fid"`
	Loaded     int64  `json:"loaded,omitempty"`
	ChunkSize  int64  `json:"chunkSize,omitempty"`
	Chunk      int    `json:"chunk,omitempty"`
	Chunks     int    `json:"chunks,omitempty"`
	Collection string `json:"collection,omitempty"`
	ChunkList  seaweedfs.ChunkList
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
