package seaweedfs

import (
	"github.com/ginuerzh/weedo"
	"net/url"
	"io"
	"github.com/chrislusf/seaweedfs/weed/operation"
)

var Client *weedo.Client

type ChunkInfo operation.ChunkInfo

type ChunkList []*ChunkInfo

type ChunkManifest struct {
	Name   string    `json:"name,omitempty"`
	Mime   string    `json:"mime,omitempty"`
	Size   int64     `json:"size,omitempty"`
	Chunks ChunkList `json:"chunks,omitempty"`
}

func GetKey(collection string) (string, error) {
	args := url.Values{}
	args.Set("collection", collection)
	return Client.Master().AssignArgs(args)
}

func Upload(fid, collection, filename, mimeType string, chunk int, file io.Reader, args url.Values) (int64, error) {
	vol, err := Client.Volume(fid, collection)
	if err != nil {
		return 0, err
	}
	return vol.Upload(fid, chunk, filename, mimeType, file, args)
}
