//go:build linux || darwin
// +build linux darwin

package memory

import (
	"os/exec"
	"regexp"
	"strings"
)

func getMemoryInfo() (memoryInfo map[string]string, err error) {
	memoryInfo = make(map[string]string)

	out, err := exec.Command("sysctl", "-n", "hw.memsize").Output()
	if err == nil {
		memoryInfo["total"] = strings.Trim(string(out), "\n")
	}

	out, err = exec.Command("sysctl", "-n", "vm.swapusage").Output()
	if err == nil {
		swap := regexp.MustCompile("total = ").Split(string(out), 2)[1]
		memoryInfo["swap_total"] = strings.Split(swap, " ")[0]
	}

	return
}

func getMemoryInfoByte() (mem uint64, swap uint64, err error) {
	memInfo, err := getMemoryInfo()
	var mem, swap uint64

	// mem is already in bytes but `swap_total` use the format "5120,00M"
	if v, ok := memInfo["swap_total"]; ok {
		idx := strings.IndexAny(v, ",.") // depending on the local either a comma or dot is used
		swapTotal, e := strconv.ParseUint(v[0:idx])
		if e == nil {
			swap = swapTotal * 1024 * 1024 // swapTotal is in mb
		}
	}

	if v, ok := memInfo["total"]; ok {
		t, e := strconv.ParseUint(v)
		if e == nil {
			mem = t // mem is returned in bytes
		}
	}

	return mem, swap, err
}
