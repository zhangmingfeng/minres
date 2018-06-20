package main

import (
	"github.com/mholt/caddy/caddy/caddymain"
	_ "github.com/hacdias/caddy-service"
	_ "github.com/zhangmingfeng/minres/plugins/router"
	_ "github.com/zhangmingfeng/minres/controllers"
)

func main() {
	caddymain.Run()
}
