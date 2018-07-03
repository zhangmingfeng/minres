package download

import (
	"github.com/zhangmingfeng/minres/controllers/base"
	"github.com/zhangmingfeng/minres/plugins/router"
	"net/http"
)

var Controller = &Download{}

func init() {
	router.RegisterController("download.image", Controller.Image)
}

type Download struct {
	base.ControllerBase
}

func (u *Download) Image(w http.ResponseWriter, r *http.Request) {
	
}
