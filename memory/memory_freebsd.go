package memory

import (
	"os/exec"
	"strings"
)

func getMemoryInfo() (memoryInfo map[string]string, err error) {
	memoryInfo = make(map[string]string)

	out, err := exec.Command("sysctl", "-n", "hw.physmem").Output()
	if err != nil {
		return memoryInfo, err
	}
	memoryInfo["total"] = strings.Trim(string(out), "\n")

	out, err = exec.Command("sysctl", "-n", "vm.swap_total").Output()
	if err != nil {
		return memoryInfo, err
	}
	memoryInfo["swap_total"] = strings.Trim(string(out), "\n")

	return
}
