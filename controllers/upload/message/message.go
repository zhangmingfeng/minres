package message

import (
	"github.com/zhangmingfeng/minres/controllers/base/message"
	"github.com/zhangmingfeng/minres/plugins/minres/weed"
)

const (
	FileNameIsEmpty    = 10001
	FileSizeIsInvalid  = 10002
	FileTimeIsInvalid  = 10003
	ChunkSizeIsInvalid = 10004
	TokenIsEmpty       = 10005
	TokenIsInvalid     = 10006
	ChunkIsInvalid     = 10007
	RemoteUrlIsEmpty   = 10008
	RemoteUrlIsInvalid = 10009
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
	Token      string `json:"token,omitempty"`
	FileHandle string `json:"fileHandle,omitempty"`
	Chunk      int    `json:"chunk,omitempty"`
}

type UploadResponse struct {
	message.BaseResponse
	IsFinished bool  `json:"isFinished,omitempty"`
	Loaded     int64 `json:"loaded,omitempty"`
	Chunk      int   `json:"chunk,omitempty"`
	Chunks     int   `json:"chunks,omitempty"`
	File       File  `json:"file,omitempty"`
}

type RemoteRequest struct {
	Url       string `json:"url,omitempty"`
	FileName  string `json:"fileName,omitempty"`
	FileGroup string `json:"fileGroup,omitempty"`
}

type RemoteResponse struct {
	message.BaseResponse
	File File `json:"file,omitempty"`
}

type File struct {
	Fid  string `json:"fid,omitempty"`
	Name string `json:"name,omitempty"`
	Size int64  `json:"size,omitempty"`
	Url  string `json:"url,omitempty"`
}

type TokenData struct {
	FileName   string            `json:"fileName,omitempty"`
	FileSize   int64             `json:"fileSize,omitempty"`
	FileTime   int64             `json:"fileTime,omitempty"`
	IsFinish   bool              `json:"isFinish,omitempty"`
	Loaded     int64             `json:"loaded,omitempty"`
	ChunkSize  int64             `json:"chunkSize,omitempty"`
	Chunk      int               `json:"chunk,omitempty"`
	Chunks     int               `json:"chunks,omitempty"`
	Collection string            `json:"collection,omitempty"`
	ChunkList  []*weed.ChunkInfo `json:"chunkList"`
}

func NewParamsRequest() *ParamsRequest {
	return &ParamsRequest{}
}

func NewUploadRequest() *UploadRequest {
	return &UploadRequest{}
}

func NewRemoteRequest() *RemoteRequest {
	return &RemoteRequest{}
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

func NewRemoteResponse() *RemoteResponse {
	return &RemoteResponse{
		BaseResponse: message.BaseResponse{
			Code: 200,
			Msg:  "success",
		},
	}
}
