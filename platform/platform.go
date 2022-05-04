//go:build !android
// +build !android

package platform

import (
	"fmt"
	"os/exec"
	"regexp"
	"runtime"
	"strings"

	"github.com/DataDog/gohai/utils"
)

// Platform holds metadata about the host
type Platform struct {
	// GoVersion is the golang version.
	GoVersion string
	// PythonVersion is the version of python in the current env (ie: returned by "python -V").
	PythonVersion string
	// GoOS is equal to "runtime.GOOS"
	GoOS string
	// GoArch is equal to "runtime.GOARCH"
	GoArch string

	// KernelName is the kernel name (ex:  "windows", "Linux", ...)
	KernelName string
	// KernelRelease the kernel release (ex: "10.0.20348", "4.15.0-1080-gcp", ...)
	KernelRelease string
	// Hostname is the hostname for the host
	Hostname string
	// Machine the architecture for the host (is: x86_64 vs arm).
	Machine string
	// OS is the os name description (ex: "GNU/Linux", "Windows Server 2022 Datacenter", ...)
	OS string

	// Family is the OS family (Windows only)
	Family string

	// KernelVersion the kernel version, Unix only
	KernelVersion string
	// Processor is the processor type, Unix only (ex "x86_64", "arm", ...)
	Processor string
	// HardwarePlatform is the hardware name, Linux only (ex "x86_64")
	HardwarePlatform string
}

const name = "platform"

func (self *Platform) Name() string {
	return name
}

func (self *Platform) Collect() (result interface{}, err error) {
	result, err = getPlatformInfo()
	return
}

// Get returns a Platform struct already initialized
func Get() (*Platform, error) {
	platformInfo, err := getPlatformInfo()
	if err != nil {
		return nil, err
	}

	p := &Platform{}
	p.GoVersion = utils.GetString(platformInfo, "goV")
	p.PythonVersion = utils.GetString(platformInfo, "pythonV")
	p.GoOS = utils.GetString(platformInfo, "GOOS")
	p.GoArch = utils.GetString(platformInfo, "GOOARCH")
	p.KernelName = utils.GetString(platformInfo, "kernel_name")
	p.KernelRelease = utils.GetString(platformInfo, "kernel_release")
	p.Hostname = utils.GetString(platformInfo, "hostname")
	p.Machine = utils.GetString(platformInfo, "machine")
	p.OS = utils.GetString(platformInfo, "os")
	p.Family = utils.GetString(platformInfo, "family")
	p.KernelVersion = utils.GetString(platformInfo, "kernel_version")
	p.Processor = utils.GetString(platformInfo, "processor")
	p.HardwarePlatform = utils.GetString(platformInfo, "hardware_platform")

	return p, nil
}

func getPlatformInfo() (platformInfo map[string]string, err error) {

	// collect each portion, and allow the parts that succeed (even if some
	// parts fail.)  For this check, it does have the (small) liability
	// that if both the ArchInfo() and the PythonVersion() fail, the error
	// from the ArchInfo() will be lost

	// for this, no error check.  The successful results will be added
	// to the return value, and the error stored.
	platformInfo, err = GetArchInfo()
	if platformInfo == nil {
		platformInfo = map[string]string{}
	}

	platformInfo["goV"] = strings.Replace(runtime.Version(), "go", "", -1)
	// If this errors, swallow the error.
	// It will usually mean that Python is not on the PATH
	// and we don't care about that.
	pythonV, e := getPythonVersion(exec.Command)

	// if there was no failure, add the python variables to the platformInfo
	if e == nil {
		platformInfo["pythonV"] = pythonV
	}

	platformInfo["GOOS"] = runtime.GOOS
	platformInfo["GOOARCH"] = runtime.GOARCH

	return
}

func getPythonVersion(execCmd utils.ExecCmdFunc) (string, error) {
	out, err := execCmd("python", "-V").CombinedOutput()
	if err != nil {
		return "", err
	}
	return parsePythonVersion(out)
}

func parsePythonVersion(cmdOut []byte) (string, error) {
	version := fmt.Sprintf("%s", cmdOut)
	values := regexp.MustCompile("Python (.*)\n").FindStringSubmatch(version)
	if len(values) < 2 {
		return "", fmt.Errorf("could not find Python version in `python -V` output: %q", version)
	}
	return strings.Trim(values[1], "\r"), nil
}
