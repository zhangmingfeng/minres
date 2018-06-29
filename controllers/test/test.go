package test

import (
	"github.com/zhangmingfeng/minres/plugins/router"
	"github.com/zhangmingfeng/minres/controllers/base"
	"net/http"
	"github.com/zhangmingfeng/minres/plugins/redis"
	"fmt"
)

var UserController = User{}

type TestRequesst struct {
	Name   string `json:"name"`
	Age    int    `json:"age"`
	Test   string `json:"test"`
	TestId int    `json:"testId"`
	Params Params `json:"params"`
}

type Params struct {
	p1 string `json:"p1"`
	p2 int    `json:"p2"`
}

func init() {
	router.RegisterController("test.Test", UserController.Test)
}

type User struct {
	base.ControllerBase
}

func (u *User) Test(w http.ResponseWriter, r *http.Request) {
	testRequest := &TestRequesst{
		Test: "aaaaaa",
	}
	u.ParseForm(r, testRequest)
	fmt.Println(testRequest)
	err := redis.Client.Set("test", map[string]string{"name": "zhangmingfeng", "age": "28"}, 0).Err()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	val, err := redis.Client.Get("test").Result()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println(val)
	r.ParseForm()
	w.WriteHeader(200)
	w.Write([]byte("hello world, test: "))
}
