package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
)

func Md5(data interface{}) (string, error) {
	dj, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	md5Ctx := md5.New()
	md5Ctx.Write(dj)
	return hex.EncodeToString(md5Ctx.Sum(nil)), nil
}
