package test

import (
	"github.com/zhangmingfeng/minres/plugins/router"
	"github.com/zhangmingfeng/minres/controllers/base"
	"net/http"
	"github.com/zhangmingfeng/minres/plugins/redis"
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
	err := redis.Client.Set("test", "ddddd", 0).Err()
	if err != nil {
		panic(err)
	}
	val, err := redis.Client.Get("test").Result()
	if err != nil {
		panic(err)
	}
	r.ParseForm()
	w.WriteHeader(200)
	w.Write([]byte("hello world, test: " + val))
}
