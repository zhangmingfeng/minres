package upload

import (
	"encoding/json"
	"fmt"
	"github.com/zhangmingfeng/minres/controllers/base"
	"github.com/zhangmingfeng/minres/controllers/upload/message"
	"github.com/zhangmingfeng/minres/plugins/minres"
	"github.com/zhangmingfeng/minres/plugins/minres/weed"
	"github.com/zhangmingfeng/minres/utils"
	"math"
	"net/http"
	"net/url"
	"runtime/debug"
	"strconv"
)

var Controller = &Upload{}

func init() {
	minres.RegisterController("/params", Controller.Params)
	minres.RegisterController("/upload", Controller.Upload)
}

type Upload struct {
	base.ControllerBase
}

/**
获取上传参数
*/
func (u *Upload) Params(w http.ResponseWriter, r *http.Request) {
	request := message.NewParamsRequest()
	response := message.NewParamsResponse()
	defer func() {
		if err := recover(); err != nil {
			response.Code = 500
			response.Msg = fmt.Sprint(err)
			u.JsonResponse(w, response)
		}
	}()
	u.ParseForm(r, request)
	if len(request.FileName) <= 0 {
		response.Code = message.FileNameIsEmpty
		response.Msg = "fileName is empty"
		u.JsonResponse(w, response)
		return
	}
	if len(request.FileGroup) <= 0 {
		request.FileGroup = "default"
	}
	if request.FileSize <= 0 {
		response.Code = message.FileSizeIsInvalid
		response.Msg = "fileSize is invalid"
		u.JsonResponse(w, response)
		return
	}
	if request.FileTime <= 0 {
		response.Code = message.FileTimeIsInvalid
		response.Msg = "fileTime is invalid"
		u.JsonResponse(w, response)
		return
	}
	if request.ChunkSize <= 0 {
		response.Code = message.ChunkSizeIsInvalid
		response.Msg = "chunkSize is invalid"
		u.JsonResponse(w, response)
		return
	}
	token, err := utils.Md5(request)
	if err != nil {
		panic(err.Error())
	}
	tokenDataBin := message.TokenData{}
	tokenData := u.CacheValue(token)
	isNew := true
	if len(tokenData) > 0 {
		err := utils.Map2Struct(tokenData, &tokenDataBin)
		if err != nil {
			panic(err.Error())
		}
		isNew = tokenDataBin.IsFinish
	}
	if isNew {
		chunks := math.Ceil(float64(request.FileSize) / float64(request.ChunkSize))
		tokenDataBin = message.TokenData{
			FileSize:   request.FileSize,
			FileName:   request.FileName,
			FileTime:   request.FileTime,
			IsFinish:   false,
			Loaded:     0,
			ChunkSize:  request.ChunkSize,
			Chunk:      1,
			Chunks:     int(chunks),
			Collection: request.FileGroup,
		}
		val, err := json.Marshal(tokenDataBin)
		if err != nil {
			panic(err.Error())
		}
		u.Cache(token, string(val), 0)
	}
	response.Token = token
	response.UploadUrl = minres.Url("/upload")
	response.ChunkSize = tokenDataBin.ChunkSize
	response.Chunk = tokenDataBin.Chunk
	response.Chunks = tokenDataBin.Chunks
	response.Loaded = tokenDataBin.Loaded
	response.FileGroup = tokenDataBin.Collection
	u.JsonResponse(w, response)
}

func (u *Upload) Upload(w http.ResponseWriter, r *http.Request) {
	request := message.NewUploadRequest()
	response := message.NewUploadResponse()
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(string(debug.Stack()))
			response.Code = 500
			response.Msg = fmt.Sprint(err)
			u.JsonResponse(w, response)
		}
	}()
	u.ParseForm(r, request)
	token := request.Token
	if len(token) <= 0 {
		response.Code = message.TokenIsEmpty
		response.Msg = "token is empty"
		u.JsonResponse(w, response)
		return
	}
	tokenData := u.CacheValue(token)
	if len(tokenData) <= 0 {
		response.Code = message.TokenIsInvalid
		response.Msg = "token is invalid"
		u.JsonResponse(w, response)
		return
	}
	tokenDataBin := message.TokenData{}
	err := utils.Map2Struct(tokenData, &tokenDataBin)
	if err != nil {
		panic(err.Error())
	}
	if request.Chunk != tokenDataBin.Chunk {
		response.Code = message.ChunkISInvalid
		response.Msg = "chunk is invalid, current chunk is " + strconv.Itoa(tokenDataBin.Chunk)
		u.JsonResponse(w, response)
		return
	}
	fileHandle := request.FileHandle
	if len(fileHandle) <= 0 {
		fileHandle = "file"
	}
	file, fileInfo, err := r.FormFile(fileHandle)
	if err != nil {
		panic(err)
	}
	args := url.Values{}
	args.Set("collection", tokenDataBin.Collection)
	fid, size, err := minres.WeedClient.Upload(fileInfo.Filename, file, args)
	if err != nil {
		panic(err)
	}
	chunkInfo := &weed.ChunkInfo{
		Fid:    fid,
		Offset: int64(request.Chunk-1) * tokenDataBin.ChunkSize,
		Size:   size,
	}
	tokenDataBin.ChunkList = append(tokenDataBin.ChunkList, chunkInfo)
	tokenDataBin.Loaded = tokenDataBin.Loaded + size
	response.Loaded = tokenDataBin.Loaded
	response.Chunk = request.Chunk
	response.Chunks = tokenDataBin.Chunks
	response.IsFinished = false
	if request.Chunk >= tokenDataBin.Chunks {
		response.IsFinished = true
		tokenDataBin.IsFinish = true
		//当总分片大于1的时候，才使用chunk功能
		if tokenDataBin.Chunks > 1 {
			chunkManifest := &weed.ChunkManifest{
				Name:   fileInfo.Filename,
				Mime:   "application/octet-stream",
				Size:   response.Loaded,
				Chunks: tokenDataBin.ChunkList,
			}
			args := url.Values{}
			args.Set("collection", tokenDataBin.Collection)
			fid, _, err = minres.WeedClient.MergeChunks(fileInfo.Filename, chunkManifest, args)
			if err != nil {
				panic(err)
			}
		}
		response.File = message.File{
			Fid: fid,
			Url: minres.Url("/fetch/{fid}", "fid", fid),
		}
	} else {
		tokenDataBin.Chunk = request.Chunk + 1
	}
	val, err := json.Marshal(tokenDataBin)
	if err != nil {
		panic(err.Error())
	}
	u.Cache(token, string(val), 0)
	u.JsonResponse(w, response)
}
