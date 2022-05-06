package cpu

import (
	"strconv"
	"strings"

	"github.com/DataDog/gohai/utils"
)

// Cpu holds metadata about the host CPU
type Cpu struct {
	// VendorId the CPU vendor ID
	VendorId string
	// ModelName the CPU model
	ModelName string
	// CpuCores the number of cores for the CPU
	CpuCores uint64
	// CpuLogicalProcessors the number of logical core for the CPU
	CpuLogicalProcessors uint64
	// Mhz the frequency for the CPU
	Mhz float64
	// CacheSize the cache size for the CPU (Linux only)
	CacheSize uint64
	// Family the CPU family
	Family string
	// Model the CPU model name
	Model string
	// Stepping the CPU stepping
	Stepping string

	// CpuPkgs the CPU pkg count (Windows only)
	CpuPkgs uint64
	// CpuNumaNodes the CPU numa node count (Windows only)
	CpuNumaNodes uint64
	// CacheSizeL1 the CPU L1 cache size (Windows only)
	CacheSizeL1 uint64
	// CacheSizeL2 the CPU L2 cache size (Windows only)
	CacheSizeL2 uint64
	// CacheSizeL3 the CPU L3 cache size (Windows only)
	CacheSizeL3 uint64
}

const name = "cpu"

func (self *Cpu) Name() string {
	return name
}

func (self *Cpu) Collect() (result interface{}, err error) {
	result, err = getCpuInfo()
	return
}

// Get returns a Cpu struct  already initialized
func Get() (*Cpu, error) {
	cpuInfo, err := getCpuInfo()
	if err != nil {
		return nil, err
	}

	c := &Cpu{}

	c.VendorId = utils.GetString(cpuInfo, "vendor_id")
	c.ModelName = utils.GetString(cpuInfo, "model_name")
	c.Family = utils.GetString(cpuInfo, "family")
	c.Model = utils.GetString(cpuInfo, "model")
	c.Stepping = utils.GetString(cpuInfo, "stepping")

	// We serialize int to string in the windows version of 'GetCpuInfo' and back to in here. This is less than
	// ideal but we don't want to break backward compatibility for now. The entire gohai project needs a rework but
	// for now we simply adding typed field to avoid using maps of interface..
	c.CpuPkgs = utils.GetUint64(cpuInfo, "cpu_pkgs")
	c.CpuNumaNodes = utils.GetUint64(cpuInfo, "cpu_numa_nodes")
	c.CacheSizeL1 = utils.GetUint64(cpuInfo, "cache_size_l1")
	c.CacheSizeL2 = utils.GetUint64(cpuInfo, "cache_size_l2")
	c.CacheSizeL3 = utils.GetUint64(cpuInfo, "cache_size_l3")

	c.CpuCores = utils.GetUint64(cpuInfo, "cpu_cores")
	c.CpuLogicalProcessors = utils.GetUint64(cpuInfo, "cpu_logical_processors")
	c.Mhz = utils.GetFloat64(cpuInfo, "mhz")

	// cache_size uses the format '9216 KB'
	cacheSizeString := strings.Split(utils.GetString(cpuInfo, "cache_size"), " ")[0]
	cacheSize, err := strconv.ParseUint(cacheSizeString, 10, 64)
	if err == nil {
		c.CacheSize = cacheSize * 1024
	}

	return c, nil
}
