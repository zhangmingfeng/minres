package upload

import (
	"github.com/zhangmingfeng/minres/controllers/base"
	"net/http"
	"github.com/zhangmingfeng/minres/controllers/upload/message"
	"fmt"
	"github.com/zhangmingfeng/minres/utils"
	"github.com/devfeel/mapper"
	"encoding/json"
	"github.com/zhangmingfeng/minres/plugins/router"
	"github.com/zhangmingfeng/minres/plugins/seaweedfs"
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
	params := message.Params{}
	tokenData := u.CacheValue(token)
	if len(tokenData) > 0 {
		tokenMap := make(map[string]interface{}, 0)
		err = json.Unmarshal([]byte(tokenData), &tokenMap)
		if err != nil {
			panic(err.Error())
		}
		mapper.MapperMap(tokenMap, &params)
	} else {
		key, err := seaweedfs.GetKey()
		if err != nil {
			panic(err.Error())
		}
		params.Key = key
		val, err := json.Marshal(params)
		if err != nil {
			panic(err.Error())
		}
		u.Cache(token, string(val), 0)
	}
	response.Params = params
	response.Token = token
	response.UploadUrl = router.Url("upload.upload")
	u.JsonResponse(w, response)
}

func (u *Upload) Upload(w http.ResponseWriter, r *http.Request) {
	request := message.NewUploadRequest()
	response := message.NewUploadResponse()
	defer func() {
		if err := recover(); err != nil {
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
}
