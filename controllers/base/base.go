package base

import (
	"net/http"
	"regexp"
	"io"
	"compress/gzip"
	"io/ioutil"
	"encoding/json"
	"fmt"
)

var (
	requestJSONRegex = regexp.MustCompile(`(application/json)(?:,|$)`)
	acceptsJSONRegex = regexp.MustCompile(`(application/json)(?:,|$)`)
	isGZIPRegex      = regexp.MustCompile("gzip")
)

type ControllerBase struct {
	request *http.Request
}

func (c *ControllerBase) ParseForm(r *http.Request) {
	c.request = r
	if c.RequestJSON() {
		bodyMap := make(map[string]string, 0)
		body := c.RequestBody()
		err := json.Unmarshal(body, &bodyMap)
		fmt.Println(bodyMap, err)
	}
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
	var requestbody []byte
	safe := &io.LimitedReader{R: c.request.Body, N: 1 << 24} //16M
	if c.IsGZIP() {
		reader, err := gzip.NewReader(safe)
		if err != nil {
			return nil
		}
		requestbody, _ = ioutil.ReadAll(reader)
	} else {
		requestbody, _ = ioutil.ReadAll(safe)
	}
	return requestbody
}
