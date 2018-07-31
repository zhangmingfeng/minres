package base

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/zhangmingfeng/minres/plugins/minres"
	"github.com/zhangmingfeng/minres/plugins/minres/log"
	"github.com/zhangmingfeng/minres/plugins/minres/weed"
	"github.com/zhangmingfeng/minres/plugins/redis"
	"github.com/zhangmingfeng/minres/utils"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	requestJSONRegex = regexp.MustCompile(`^(application/json)`)
	acceptsJSONRegex = regexp.MustCompile(`^(application/json)`)
	isGZIPRegex      = regexp.MustCompile("gzip")
)

type ControllerBase struct {
	request *http.Request
}

func (c *ControllerBase) ParseForm(r *http.Request, request interface{}) {
	c.request = r
	paramsMap := make(map[string]interface{}, 0)
	vars := mux.Vars(r)
	for k, v := range vars {
		paramsMap[k] = v
	}
	if c.RequestJSON() {
		bodyMap := make(map[string]interface{}, 0)
		body := c.RequestBody()
		err := json.Unmarshal(body, &bodyMap)
		if err != nil {
			panic(err)
		}
		for k, v := range bodyMap {
			paramsMap[k] = v
		}
	}
	r.ParseMultipartForm(32 << 20)
	for k, _ := range r.Form {
		paramsMap[k] = r.FormValue(k)
	}
	utils.Map2Struct(paramsMap, request)
}

func (c *ControllerBase) RequestJSON() bool {
	return requestJSONRegex.MatchString(c.request.Header.Get("Content-Type"))
}

func (c *ControllerBase) AcceptsJSON() bool {
	return acceptsJSONRegex.MatchString(c.request.Header.Get("Accept"))
}

func (c *ControllerBase) IsGZIP() bool {
	return isGZIPRegex.MatchString(c.request.Header.Get("Content-Encoding"))
}

func (c *ControllerBase) RequestBody() []byte {
	var requestBody []byte
	safe := &io.LimitedReader{R: c.request.Body, N: 1 << 24} //16M
	if c.IsGZIP() {
		reader, err := gzip.NewReader(safe)
		if err != nil {
			return nil
		}
		requestBody, _ = ioutil.ReadAll(reader)
	} else {
		requestBody, _ = ioutil.ReadAll(safe)
	}
	return requestBody
}

func (c *ControllerBase) Cache(key string, val string, ttl time.Duration) error {
	err := redis.Set(key, val, ttl)
	if err != nil {
		panic(err.Error())
	}
	return nil
}

func (c *ControllerBase) CacheValue(key string) string {
	val, err := redis.Get(key)
	if err != nil && err != redis.Nil {
		panic(err)
	}
	return val
}

func (c *ControllerBase) JsonResponse(w http.ResponseWriter, res interface{}) {
	body, _ := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (c *ControllerBase) TextResponse(w http.ResponseWriter, status int, text string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	w.Write([]byte(text))
}

func (c *ControllerBase) FileResponse(w http.ResponseWriter, r *http.Request, fileInfo *weed.FileInfo, isDownload bool) error {
	if fileInfo.Mime == "" {
		if ext := path.Ext(fileInfo.Name); ext != "" {
			fileInfo.Mime = mime.TypeByExtension(ext)
		}
	}
	w.Header().Set("Content-Type", fileInfo.Mime)
	contentDisposition := "inline"
	if !strings.HasPrefix(fileInfo.Mime, "image") {
		contentDisposition = "attachment"
	}
	if isDownload {
		contentDisposition = "attachment"
	}
	w.Header().Set("Content-Disposition", contentDisposition+`; filename="`+fileInfo.Name+`"`)
	w.Header().Set("Content-Length", strconv.FormatInt(fileInfo.Size, 10))
	if r.Method == "HEAD" {
		return nil
	}
	_, err := io.Copy(w, bytes.NewReader(fileInfo.GetData()))
	return err
}

func (c *ControllerBase) Logger() *log.Logger {
	return minres.Logger
}
