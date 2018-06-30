package seaweedfs

import (
	"github.com/ginuerzh/weedo"
)

var Client *weedo.Client

func GetKey() (string, error) {
	return Client.Master().Assign()
}
