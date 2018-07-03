package test

import (
	"github.com/zhangmingfeng/minres/plugins/router"
	"github.com/zhangmingfeng/minres/controllers/base"
	"net/http"
	"github.com/zhangmingfeng/minres/plugins/seaweedfs"
	"fmt"
)

var UserController = User{}

type TestRequesst struct {
	Name   string   `json:"name"`
	Age    int      `json:"age"`
	Fid    string   `json:"fid"`
	Params []*Param `json:"params"`
}

type Param struct {
	P1 string `json:"p1"`
	P2 int    `json:"p2"`
}

func init() {
	router.RegisterController("test.Test", UserController.Test)
}

type User struct {
	base.ControllerBase
}

func (u *User) Test(w http.ResponseWriter, r *http.Request) {
	request := TestRequesst{}
	u.ParseForm(r, &request)
	err := seaweedfs.Delete(request.Fid, "category")
	fmt.Println(err, request.Fid)
}
