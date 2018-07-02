package test

import (
	"github.com/zhangmingfeng/minres/plugins/router"
	"github.com/zhangmingfeng/minres/controllers/base"
	"net/http"
	"encoding/json"
	"fmt"
)

var UserController = User{}

type TestRequesst struct {
	Name   string   `json:"name"`
	Age    int      `json:"age"`
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
	testR := &TestRequesst{
		Name: "zhang",
		Age:  28,
		Params: []*Param{
			&Param{
				P1: "1p1",
				P2: 1,
			},
			&Param{
				P1: "2p1",
				P2: 3,
			}},
	}
	fmt.Println(testR, testR.Params[0], testR.Params[1])
	b, _ := json.Marshal(testR)
	testR2 := &TestRequesst{}
	err := json.Unmarshal(b, testR2)
	fmt.Println(err)
	fmt.Println(testR2, testR2.Params[0], testR2.Params[1])
}
