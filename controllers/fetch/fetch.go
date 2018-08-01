package fetch

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zhangmingfeng/minres/controllers/base"
	"github.com/zhangmingfeng/minres/controllers/fetch/message"
	"github.com/zhangmingfeng/minres/plugins/minres"
	"github.com/zhangmingfeng/minres/plugins/minres/weed"
	"github.com/zhangmingfeng/minres/utils"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var Controller = &Fetch{}

func init() {
	minres.RegisterController("/fetch/{fid}", Controller.Fetch)
}

type Fetch struct {
	base.ControllerBase
}

func (f *Fetch) Fetch(w http.ResponseWriter, r *http.Request) {
	request := message.Request{}
	f.ParseForm(r, &request)
	defer func() {
		if err := recover(); err != nil {
			f.Logger().Error(err)
			f.TextResponse(w, http.StatusNotFound, "not found")
		}
	}()
	fid := request.Fid
	if len(fid) == 0 {
		panic(errors.New("not found"))
	}
	width := request.Width
	height := request.Height
	mode := request.Mode
	prefixList := make([]string, 0)
	args := url.Values{}
	if width > 0 {
		args.Set("width", strconv.Itoa(width))
		prefixList = append(prefixList, fmt.Sprintf("w%d", width))
	}
	if height > 0 {
		args.Set("height", strconv.Itoa(height))
		prefixList = append(prefixList, fmt.Sprintf("h%d", height))
	}
	if len(mode) > 0 {
		args.Set("mode", mode)
		prefixList = append(prefixList, fmt.Sprintf("m%s", mode))
	}
	prefix := strings.Join(prefixList, "_")
	var fileMeta string
	if len(prefix) > 0 {
		fileMeta = f.CacheValue(fmt.Sprintf("%s_%s_meta", fid, prefix))
	}
	var fileInfo *weed.FileInfo
	var err error
	if len(fileMeta) <= 0 {
		fileInfo, err = minres.WeedClient.Fetch(fid, args)
		if err != nil {
			panic(err)
		}
		// 仅仅缩略图才缓存文件
		if strings.HasPrefix(fileInfo.Mime, "image") && len(prefix) > 0 {
			err = minres.WeedClient.SaveFile(fmt.Sprintf("%s_%s", fid, prefix), fileInfo.GetData())
			if err != nil {
				panic(err)
			}
			val, err := json.Marshal(fileInfo)
			if err != nil {
				panic(err)
			}
			f.Cache(fmt.Sprintf("%s_%s_meta", fid, prefix), string(val), 0)
			f.Logger().Debug("save cache", val)
		}
	} else {
		f.Logger().Debug("from cache", fileMeta)
		fileInfo = &weed.FileInfo{}
		err := utils.Json2Struct(fileMeta, fileInfo)
		if err != nil {
			panic(err)
		}
		data, err := minres.WeedClient.ReadFile(fmt.Sprintf("%s_%s", fid, prefix))
		if err != nil { //如果读文件失败，就重新获取
			fileInfo, err := minres.WeedClient.Fetch(fid, args)
			if err != nil {
				panic(err)
			}
			err = minres.WeedClient.SaveFile(fmt.Sprintf("%s_%s", fid, prefix), fileInfo.GetData())
			if err != nil {
				panic(err)
			}
		}
		fileInfo.SetData(data)
	}
	err = f.FileResponse(w, r, fileInfo, request.Download)
	if err != nil {
		panic(err)
	}
}
