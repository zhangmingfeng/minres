package upload

import (
	"encoding/json"
	"fmt"
	"github.com/zhangmingfeng/minres/controllers/base"
	"github.com/zhangmingfeng/minres/controllers/upload/message"
	"github.com/zhangmingfeng/minres/plugins/router"
	"github.com/zhangmingfeng/minres/plugins/seaweedfs"
	"github.com/zhangmingfeng/minres/utils"
	"math"
	"net/http"
	"runtime/debug"
	"net/url"
	"bytes"
)

var Controller = &Upload{}

func init() {
	router.RegisterController("upload.params", Controller.Params)
	router.RegisterController("upload.upload", Controller.Upload)
}

type Upload struct {
	base.ControllerBase
}

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
	if len(tokenData) > 0 {
		err := utils.Map2Struct(tokenData, &tokenDataBin)
		if err != nil {
			panic(err.Error())
		}
	} else {
		fid, err := seaweedfs.GetKey(request.FileGroup)
		if err != nil {
			panic(err.Error())
		}
		tokenDataBin.Fid = fid
		tokenDataBin.Loaded = 0
		tokenDataBin.ChunkSize = request.ChunkSize
		tokenDataBin.Chunk = 1
		chunks := math.Ceil(float64(request.FileSize) / float64(request.ChunkSize))
		tokenDataBin.Chunks = int(chunks)
		tokenDataBin.Collection = request.FileGroup
		val, err := json.Marshal(tokenDataBin)
		if err != nil {
			panic(err.Error())
		}
		u.Cache(token, string(val), 0)
	}
	response.Token = token
	response.UploadUrl = router.Url("upload.upload")
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
	file, fileInfo, err := r.FormFile("file")
	if err != nil {
		panic(err)
	}
	fid, err := seaweedfs.GetKey("category")
	fmt.Println(fid)
	if err != nil {
		panic(err.Error())
	}
	size, err := seaweedfs.Upload(fid, tokenDataBin.Collection, fileInfo.Filename, "", 0, file, url.Values{})
	if err != nil {
		panic(err)
	}
	tokenDataBin.Loaded = tokenDataBin.Loaded + size
	chunkInfo := &seaweedfs.ChunkInfo{}
	chunkInfo.Fid = fid
	chunkInfo.Offset = 0
	chunkInfo.Size = size
	tokenDataBin.ChunkList = append(tokenDataBin.ChunkList, chunkInfo)
	val, err := json.Marshal(tokenDataBin)
	if err != nil {
		panic(err.Error())
	}
	u.Cache(token, string(val), 0)
	response.Loaded = tokenDataBin.Loaded
	response.Chunk = request.Chunk
	response.Chunks = tokenDataBin.Chunks
	response.IsFinished = false
	if request.Chunk >= tokenDataBin.Chunks {
		response.IsFinished = true
		chunkManifest := &seaweedfs.ChunkManifest{
			Name:   fileInfo.Filename,
			Mime:   "application/octet-stream",
			Size:   response.Loaded,
			Chunks: tokenDataBin.ChunkList,
		}
		data, _ := json.Marshal(chunkManifest)
		fmt.Println(string(data))
		buffer := bytes.NewBuffer(data)
		args := url.Values{}
		args.Set("cm", "true")
		fid, err := seaweedfs.GetKey("category")
		if err != nil {
			panic(err.Error())
		}
		response.File = message.File{
			Fid: fid,
		}
		_, err = seaweedfs.Upload(fid, tokenDataBin.Collection, fileInfo.Filename, "application/json", 0, buffer, args)
		if err != nil {
			panic(err)
		}
	}
	u.JsonResponse(w, response)
}
