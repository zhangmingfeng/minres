package main

import (
	_ "github.com/captncraig/cors/caddy"
	_ "github.com/hacdias/caddy-service"
	"github.com/mholt/caddy/caddy/caddymain"
	_ "github.com/zhangmingfeng/minres/controllers"
	_ "github.com/zhangmingfeng/minres/plugins/redis"
	_ "github.com/zhangmingfeng/minres/plugins/minres"
)

func main() {
	caddymain.Run()
}
