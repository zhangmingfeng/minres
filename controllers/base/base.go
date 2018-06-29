package base

import (
	"net/http"
	"regexp"
	"io"
	"compress/gzip"
	"io/ioutil"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/devfeel/mapper"
	"github.com/zhangmingfeng/minres/plugins/redis"
	"time"
)

var (
	requestJSONRegex = regexp.MustCompile(`(application/json)(?:,|$)`)
	acceptsJSONRegex = regexp.MustCompile(`(application/json)(?:,|$)`)
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
	r.ParseForm()
	for k, _ := range r.Form {
		paramsMap[k] = r.FormValue(k)
	}
	mapper.MapperMap(paramsMap, request)
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
