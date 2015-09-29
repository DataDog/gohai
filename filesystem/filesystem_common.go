package filesystem

import (
	"time"
	timer "github.com/DataDog/gohai/utils"
)

type FileSystem struct{}

const name = "filesystem"

func (self *FileSystem) Name() string {
	return name
}

func (self *FileSystem) Collect() (result interface{}, err error) {
	defer timer.TimeTrack(time.Now(), "filesystem")
	result, err = getFileSystemInfo()
	return
}
