package test

import (
	"github.com/zhangmingfeng/minres/plugins/router"
	"github.com/zhangmingfeng/minres/controllers/base"
	"net/http"
)

var UserController = User{}

func init() {
	router.RegisterController("test.Test", UserController.Test)
}

type User struct {
	base.ControllerBase
}

func (u *User) Test(w http.ResponseWriter, r *http.Request) {
	u.ParseForm(r)
	r.ParseForm()
	w.WriteHeader(200)
	w.Write([]byte("hello world"))
}
