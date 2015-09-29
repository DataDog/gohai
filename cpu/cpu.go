package cpu
import (
	"time"
	timer "github.com/DataDog/gohai/utils"
)

type Cpu struct{}

const name = "cpu"

func (self *Cpu) Name() string {
	return name
}

func (self *Cpu) Collect() (result interface{}, err error) {
	defer timer.TimeTrack(time.Now(), "cpu")
	result, err = getCpuInfo()
	return
}
