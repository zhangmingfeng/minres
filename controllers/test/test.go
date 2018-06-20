package test

import (
	"github.com/zhangmingfeng/minres/plugins/router"
	"net/http"
)

func init() {
	router.RegisterController("test.Test", Test)
}

func Test(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("hello world"))
}
