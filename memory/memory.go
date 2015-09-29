package memory
import (
	"time"
	timer "github.com/DataDog/gohai/utils"
)

type Memory struct{}

const name = "memory"

func (self *Memory) Name() string {
	return name
}

func (self *Memory) Collect() (result interface{}, err error) {
	defer timer.TimeTrack(time.Now(), "memory")
	result, err = getMemoryInfo()
	return
}
