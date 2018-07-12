package fetch

import (
	"github.com/zhangmingfeng/minres/controllers/base"
	"github.com/zhangmingfeng/minres/controllers/fetch/message"
	"github.com/zhangmingfeng/minres/plugins/router"
	"net/http"
	"fmt"
	"strings"
	"github.com/zhangmingfeng/minres/plugins/seaweedfs"
	"github.com/zhangmingfeng/mapper"
	"encoding/json"
	"time"
	"errors"
	"net/url"
	"strconv"
	"runtime/debug"
)

var Controller = &Fetch{}

func init() {
	router.RegisterController("fetch.fetch", Controller.Fetch)
}

type Fetch struct {
	base.ControllerBase
}

func (f *Fetch) Fetch(w http.ResponseWriter, r *http.Request) {
	request := message.Request{}
	f.ParseForm(r, &request)
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			debug.Stack()
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
	if width > 0 {
		prefixList = append(prefixList, fmt.Sprintf("w%d", width))
	}
	if height > 0 {
		prefixList = append(prefixList, fmt.Sprintf("h%d", height))
	}
	if len(mode) > 0 {
		prefixList = append(prefixList, fmt.Sprintf("m%s", mode))
	}
	prefix := strings.Join(prefixList, "_")
	fileMeta := f.CacheValue(fmt.Sprintf("%s_%s_meta", fid, prefix))
	args := url.Values{}
	args.Set("width", strconv.Itoa(width))
	args.Set("height", strconv.Itoa(height))
	args.Set("mode", mode)
	var fileInfo *seaweedfs.FileInfo
	var err error
	if len(fileMeta) <= 0 {
		fileInfo, err = seaweedfs.Fetch(fid, args)
		if err != nil {
			panic(err)
		}
		err = seaweedfs.SaveFile(fmt.Sprintf("%s_%s", fid, prefix), fileInfo.GetData())
		if err != nil {
			panic(err)
		}
		val, err := json.Marshal(fileInfo)
		if err != nil {
			panic(err)
		}
		f.Cache(fmt.Sprintf("%s_%s_meta", fid, prefix), string(val), 3600*24*time.Second)
	} else {
		fileInfo = &seaweedfs.FileInfo{}
		err := mapper.Json2Struct(fileMeta, fileInfo)
		if err != nil {
			panic(err)
		}
		data, err := seaweedfs.ReadFile(fmt.Sprintf("%s_%s", fid, prefix))
		if err != nil { //如果读文件失败，就重新获取
			fileInfo, err := seaweedfs.Fetch(fid, args)
			if err != nil {
				panic(err)
			}
			err = seaweedfs.SaveFile(fmt.Sprintf("%s_%s", fid, prefix), fileInfo.GetData())
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
