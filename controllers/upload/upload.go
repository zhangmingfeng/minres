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
	"mime"
)

var Controller = &Upload{}

func init() {
	minres.RegisterController("/params", Controller.Params)
	minres.RegisterController("/upload", Controller.Upload)
	minres.RegisterController("/remote", Controller.Remote)
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
			u.Logger().Error(err)
			response.Code = 500
			response.Msg = fmt.Sprint(err)
			u.JsonResponse(w, response)
		}
	}()
	u.ParseForm(r, request)
	if len(request.FileName) <= 0 {
		u.Logger().Error("fileName is empty")
		response.Code = message.FileNameIsEmpty
		response.Msg = "fileName is empty"
		u.JsonResponse(w, response)
		return
	}
	if len(request.FileGroup) <= 0 {
		request.FileGroup = "default"
	}
	if request.FileSize <= 0 {
		u.Logger().Error("fileSize is invalid")
		response.Code = message.FileSizeIsInvalid
		response.Msg = "fileSize is invalid"
		u.JsonResponse(w, response)
		return
	}
	if request.FileTime <= 0 {
		u.Logger().Error("fileTime is invalid")
		response.Code = message.FileTimeIsInvalid
		response.Msg = "fileTime is invalid"
		u.JsonResponse(w, response)
		return
	}
	if request.ChunkSize <= 0 {
		u.Logger().Error("chunkSize is invalid")
		response.Code = message.ChunkSizeIsInvalid
		response.Msg = "chunkSize is invalid"
		u.JsonResponse(w, response)
		return
	}
	token, err := utils.Md5(request)
	if err != nil {
		u.Logger().Error(err.Error())
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
			u.Logger().Error(err)
			response.Code = 500
			response.Msg = fmt.Sprint(err)
			u.JsonResponse(w, response)
		}
	}()
	u.ParseForm(r, request)
	token := request.Token
	if len(token) <= 0 {
		u.Logger().Error("token is empty")
		response.Code = message.TokenIsEmpty
		response.Msg = "token is empty"
		u.JsonResponse(w, response)
		return
	}
	tokenData := u.CacheValue(token)
	if len(tokenData) <= 0 {
		u.Logger().Error("token is invalid")
		response.Code = message.TokenIsInvalid
		response.Msg = "token is invalid"
		u.JsonResponse(w, response)
		return
	}
	tokenDataBin := message.TokenData{}
	err := utils.Json2Struct(tokenData, &tokenDataBin)
	if err != nil {
		panic(err.Error())
	}
	if request.Chunk != tokenDataBin.Chunk {
		response.Code = message.ChunkIsInvalid
		response.Msg = "chunk is invalid, current chunk is " + strconv.Itoa(tokenDataBin.Chunk)
		u.Logger().Error(response.Msg)
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
			Fid:  fid,
			Name: fileInfo.Filename,
			Size: response.Loaded,
			Url:  minres.Url("/fetch/{fid}", "fid", fid),
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

func (u *Upload) Remote(w http.ResponseWriter, r *http.Request) {
	request := message.NewRemoteRequest()
	response := message.NewRemoteResponse()
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(string(debug.Stack()))
			u.Logger().Error(err)
			response.Code = 500
			response.Msg = fmt.Sprint(err)
			u.JsonResponse(w, response)
		}
	}()
	u.ParseForm(r, request)
	remoteUrl := request.Url
	if len(remoteUrl) <= 0 {
		u.Logger().Error("url is empty")
		response.Code = message.RemoteUrlIsEmpty
		response.Msg = "url is empty"
		u.JsonResponse(w, response)
		return
	}
	if len(request.FileGroup) <= 0 {
		request.FileGroup = "default"
	}
	if len(request.FileName) <= 0 {
		request.FileName, _ = utils.Md5(request)
	}
	resp, err := http.Get(remoteUrl)
	if err != nil {
		u.Logger().Error(err.Error())
		response.Code = message.RemoteUrlIsInvalid
		response.Msg = "url is invalid"
		u.JsonResponse(w, response)
		return
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		u.Logger().Error("url is invalid, httpCode: ", resp.StatusCode)
		response.Code = message.RemoteUrlIsInvalid
		response.Msg = "url is invalid"
		u.JsonResponse(w, response)
		return
	}
	fileExt := getExtByHttpHeader(resp.Header)
	if len(fileExt) == 0 {
		u.Logger().Error("remote file content-type is invalid, content-type: ", resp.Header.Get("Content-Type"))
		response.Code = message.RemoteUrlIsInvalid
		response.Msg = "url is invalid"
		u.JsonResponse(w, response)
		return
	}
	args := url.Values{}
	args.Set("collection", request.FileGroup)
	request.FileName = fmt.Sprintf("%s%s", request.FileName, fileExt)
	fid, size, err := minres.WeedClient.Upload(request.FileName, resp.Body, args)
	if err != nil {
		panic(err)
	}
	response.File = message.File{
		Fid:  fid,
		Name: request.FileName,
		Size: size,
		Url:  minres.Url("/fetch/{fid}", "fid", fid),
	}
	u.JsonResponse(w, response)
}

func getExtByHttpHeader(header http.Header) string {
	contentType := header.Get("Content-Type")
	exts, err := mime.ExtensionsByType(contentType)
	if err != nil {
		return ""
	}
	if len(exts) == 0 {
		return ""
	}
	return exts[0]
}
