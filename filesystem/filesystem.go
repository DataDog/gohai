package filesystem

import (
	"strconv"
)

type FileSystem struct{}

type MountInfo struct {
	Name      string `json:"name"`
	SizeKB    uint64 `json:"kb_size"`
	MountedOn string `json:"mounted_on"`
}

const name = "filesystem"

func (self *FileSystem) Name() string {
	return name
}

func (self *FileSystem) Collect() (interface{}, error) {
	mounts, err := self.Get()
	if err != nil {
		return nil, err
	}

	results := make([]interface{}, len(mounts))
	for idx, mount := range mounts {
		tmpMount := mount
		results[idx] = map[string]string{
			"name":       tmpMount.Name,
			"kb_size":    strconv.FormatUint(tmpMount.SizeKB, 10),
			"mounted_on": tmpMount.MountedOn,
		}
	}

	return results, nil
}

func (self *FileSystem) Get() ([]MountInfo, error) {
	return getFileSystemInfo()
}
