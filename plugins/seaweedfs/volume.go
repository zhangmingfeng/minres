package seaweedfs

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"bytes"
)

type Volume struct {
	Locations []Location
}

func NewVolume(locations []Location) *Volume {
	for i, _ := range locations {
		if !strings.HasPrefix(locations[i].Url, "http:") {
			locations[i].Url = fmt.Sprintf("http://%s", locations[i].Url)
		}
		if !strings.HasPrefix(locations[i].PublicUrl, "http:") {
			locations[i].PublicUrl = fmt.Sprintf("http://%s", locations[i].PublicUrl)
		}
	}
	return &Volume{Locations: locations}
}

func (v *Volume) Upload(fid, filename string, file io.Reader) (size int64, err error) {
	formData, contentType, err := makeFormData(filename, "application/octet-stream", file)
	if err != nil {
		return
	}
	u := fmt.Sprintf("%s/%s", v.PublicUrl(), fid)
	resp, err := upload(u, contentType, formData)
	if err == nil {
		size = resp.Size
	}

	return
}

func (v *Volume) MergeChunks(fid, filename string, chunkManifest *ChunkManifest) (err error) {
	data, err := json.Marshal(chunkManifest)
	if err != nil {
		return
	}
	formData, contentType, err := makeFormData(filename, "application/json;charset=UTF-8", bytes.NewBuffer(data))
	if err != nil {
		return
	}
	u := fmt.Sprintf("%s/%s?cm=true", v.PublicUrl(), fid)
	_, err = upload(u, contentType, formData)
	return
}

func (v *Volume) Delete(fid string, count int) (err error) {
	if count <= 0 {
		count = 1
	}

	url := fmt.Sprintf("%s/%s", v.PublicUrl(), fid)
	if err = del(url); err != nil {
		return err
	}

	for i := 1; i < count; i++ {
		if err = del(fmt.Sprintf("%s_%d", url, i)); err != nil {
			return
		}
	}
	return
}

func (v *Volume) AssignVolume(volumeId uint64, replica string) error {
	values := url.Values{}
	values.Set("volume", strconv.FormatUint(volumeId, 10))
	if len(replica) > 0 {
		values.Set("replication", replica)
	}

	_, err := http.Get(fmt.Sprintf("%s/admin/assign_volume?%s", v.PublicUrl(), values.Encode()))
	return err
}

func (v *Volume) Url() string {
	if len(v.Locations) == 0 {
		return ""
	}
	return v.Locations[0].Url
}

func (v *Volume) PublicUrl() string {
	if len(v.Locations) == 0 {
		return ""
	}
	return v.Locations[0].PublicUrl
}

type volumeStatus struct {
	Version string
	volumes []volume
	Error   string
}

type volume struct {
	Id               uint64
	Size             uint64
	RepType          string
	Version          int
	FileCount        uint64
	DeleteCount      uint64
	DeletedByteCount uint64
	ReadOnly         bool
}

// Check Volume Server Status
func (v *Volume) Status() (err error) {
	resp, err := http.Get(fmt.Sprintf("%s/status", v.PublicUrl()))
	if err != nil {
		return
	}

	defer resp.Body.Close()

	status := new(volumeStatus)
	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(status); err != nil {
		return
	}

	if status.Error != "" {
		err = errors.New(status.Error)
		return
	}
	return
}
