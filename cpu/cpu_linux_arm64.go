// +build linux
// +build arm64

package cpu

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
)

// The Linux kernel does not include much useful information in /proc/cpuinfo
// for arm64, so we must dig further into the /sys tree and build a more
// accurate representation of the contained data, rather than relying on the
// simple analysis in cpu/cpu_linux_default.go.

// nodeNRegex recognizes directories named `nodeNN`
var nodeNRegex = regexp.MustCompile("^node[0-9]+$")

// hwImpl defines values for a spcific "CPU Implementer"
type hwImpl struct {
	// name of the implementer
	name string
	// part numbers (indexed by "CPU part")
	parts map[uint64]string
}

// hwVariant is based on hw_implementer in
//     https://github.com/util-linux/util-linux/blob/master/sys-utils/lscpu-arm.c
// See from-lscpu-arm.py to regenerate this value.
var hwVariant = map[uint64]hwImpl{
	0x41: hwImpl{
		name: "ARM",
		parts: map[uint64]string{
			0x810: "ARM810",
			0x920: "ARM920",
			0x922: "ARM922",
			0x926: "ARM926",
			0x940: "ARM940",
			0x946: "ARM946",
			0x966: "ARM966",
			0xa20: "ARM1020",
			0xa22: "ARM1022",
			0xa26: "ARM1026",
			0xb02: "ARM11 MPCore",
			0xb36: "ARM1136",
			0xb56: "ARM1156",
			0xb76: "ARM1176",
			0xc05: "Cortex-A5",
			0xc07: "Cortex-A7",
			0xc08: "Cortex-A8",
			0xc09: "Cortex-A9",
			0xc0d: "Cortex-A17",
			0xc0f: "Cortex-A15",
			0xc0e: "Cortex-A17",
			0xc14: "Cortex-R4",
			0xc15: "Cortex-R5",
			0xc17: "Cortex-R7",
			0xc18: "Cortex-R8",
			0xc20: "Cortex-M0",
			0xc21: "Cortex-M1",
			0xc23: "Cortex-M3",
			0xc24: "Cortex-M4",
			0xc27: "Cortex-M7",
			0xc60: "Cortex-M0+",
			0xd01: "Cortex-A32",
			0xd03: "Cortex-A53",
			0xd04: "Cortex-A35",
			0xd05: "Cortex-A55",
			0xd06: "Cortex-A65",
			0xd07: "Cortex-A57",
			0xd08: "Cortex-A72",
			0xd09: "Cortex-A73",
			0xd0a: "Cortex-A75",
			0xd0b: "Cortex-A76",
			0xd0c: "Neoverse-N1",
			0xd0d: "Cortex-A77",
			0xd0e: "Cortex-A76AE",
			0xd13: "Cortex-R52",
			0xd20: "Cortex-M23",
			0xd21: "Cortex-M33",
			0xd40: "Neoverse-V1",
			0xd41: "Cortex-A78",
			0xd42: "Cortex-A78AE",
			0xd44: "Cortex-X1",
			0xd46: "Cortex-A510",
			0xd47: "Cortex-A710",
			0xd48: "Cortex-X2",
			0xd49: "Neoverse-N2",
			0xd4a: "Neoverse-E1",
			0xd4b: "Cortex-A78C",
			0xd4d: "Cortex-A715",
			0xd4e: "Cortex-X3",
		},
	},
	0x42: hwImpl{
		name: "Broadcom",
		parts: map[uint64]string{
			0x0f:  "Brahma B15",
			0x100: "Brahma B53",
			0x516: "ThunderX2",
		},
	},
	0x43: hwImpl{
		name: "Cavium",
		parts: map[uint64]string{
			0x0a0: "ThunderX",
			0x0a1: "ThunderX 88XX",
			0x0a2: "ThunderX 81XX",
			0x0a3: "ThunderX 83XX",
			0x0af: "ThunderX2 99xx",
		},
	},
	0x44: hwImpl{
		name: "DEC",
		parts: map[uint64]string{
			0xa10: "SA110",
			0xa11: "SA1100",
		},
	},
	0x46: hwImpl{
		name: "FUJITSU",
		parts: map[uint64]string{
			0x001: "A64FX",
		},
	},
	0x48: hwImpl{
		name: "HiSilicon",
		parts: map[uint64]string{
			0xd01: "Kunpeng-920",
		},
	},
	0x49: hwImpl{
		name:  "Infineon",
		parts: map[uint64]string{},
	},
	0x4d: hwImpl{
		name:  "Motorola/Freescale",
		parts: map[uint64]string{},
	},
	0x4e: hwImpl{
		name: "NVIDIA",
		parts: map[uint64]string{
			0x000: "Denver",
			0x003: "Denver 2",
			0x004: "Carmel",
		},
	},
	0x50: hwImpl{
		name: "APM",
		parts: map[uint64]string{
			0x000: "X-Gene",
		},
	},
	0x51: hwImpl{
		name: "Qualcomm",
		parts: map[uint64]string{
			0x00f: "Scorpion",
			0x02d: "Scorpion",
			0x04d: "Krait",
			0x06f: "Krait",
			0x201: "Kryo",
			0x205: "Kryo",
			0x211: "Kryo",
			0x800: "Falkor V1/Kryo",
			0x801: "Kryo V2",
			0x803: "Kryo 3XX Silver",
			0x804: "Kryo 4XX Gold",
			0x805: "Kryo 4XX Silver",
			0xc00: "Falkor",
			0xc01: "Saphira",
		},
	},
	0x53: hwImpl{
		name: "Samsung",
		parts: map[uint64]string{
			0x001: "exynos-m1",
		},
	},
	0x56: hwImpl{
		name: "Marvell",
		parts: map[uint64]string{
			0x131: "Feroceon 88FR131",
			0x581: "PJ4/PJ4b",
			0x584: "PJ4B-MP",
		},
	},
	0x61: hwImpl{
		name: "Apple",
		parts: map[uint64]string{
			0x020: "Icestorm-T8101",
			0x021: "Firestorm-T8101",
			0x022: "Icestorm-T8103",
			0x023: "Firestorm-T8103",
			0x030: "Blizzard-T8110",
			0x031: "Avalanche-T8110",
			0x032: "Blizzard-T8112",
			0x033: "Avalanche-T8112",
		},
	},
	0x66: hwImpl{
		name: "Faraday",
		parts: map[uint64]string{
			0x526: "FA526",
			0x626: "FA626",
		},
	},
	0x69: hwImpl{
		name: "Intel",
		parts: map[uint64]string{
			0x200: "i80200",
			0x210: "PXA250A",
			0x212: "PXA210A",
			0x242: "i80321-400",
			0x243: "i80321-600",
			0x290: "PXA250B/PXA26x",
			0x292: "PXA210B",
			0x2c2: "i80321-400-B0",
			0x2c3: "i80321-600-B0",
			0x2d0: "PXA250C/PXA255/PXA26x",
			0x2d2: "PXA210C",
			0x411: "PXA27x",
			0x41c: "IPX425-533",
			0x41d: "IPX425-400",
			0x41f: "IPX425-266",
			0x682: "PXA32x",
			0x683: "PXA930/PXA935",
			0x688: "PXA30x",
			0x689: "PXA31x",
			0xb11: "SA1110",
			0xc12: "IPX1200",
		},
	},
	0x70: hwImpl{
		name: "Phytium",
		parts: map[uint64]string{
			0x660: "FTC660",
			0x661: "FTC661",
			0x662: "FTC662",
			0x663: "FTC663",
		},
	},
	0xc0: hwImpl{
		name:  "Ampere",
		parts: map[uint64]string{},
	},
}

func getCpuInfo() (cpuInfo map[string]string, err error) {
	cpuInfo = make(map[string]string)

	procCpu, err := readProcCpuInfo()
	if err != nil {
		return nil, err
	}

	// we blithely assume that many of the CPU characteristics are the same for
	// all CPUs, so we can just use the first.
	firstCpu := procCpu[0]

	// determine vendor and model from CPU implementer / part
	if cpuVariantStr, ok := firstCpu["CPU implementer"]; ok {
		if cpuVariant, err := strconv.ParseUint(cpuVariantStr, 0, 64); err == nil {
			if cpuPartStr, ok := firstCpu["CPU part"]; ok {
				if cpuPart, err := strconv.ParseUint(cpuPartStr, 0, 64); err == nil {
					cpuInfo["model"] = cpuPartStr
					if impl, ok := hwVariant[cpuVariant]; ok {
						cpuInfo["vendor_id"] = impl.name
						if modelName, ok := impl.parts[cpuPart]; ok {
							cpuInfo["model_name"] = modelName
						} else {
							cpuInfo["model_name"] = cpuPartStr
						}
					} else {
						cpuInfo["vendor_id"] = cpuVariantStr
						cpuInfo["model_name"] = cpuPartStr
					}
				}
			}
		}
	}

	// ARM does not define a family
	cpuInfo["family"] = "none"

	// 'lscpu' represents the stepping as an rXpY string
	if cpuVariantStr, ok := firstCpu["CPU variant"]; ok {
		if cpuVariant, err := strconv.ParseUint(cpuVariantStr, 0, 64); err == nil {
			if cpuRevisionStr, ok := firstCpu["CPU revision"]; ok {
				if cpuRevision, err := strconv.ParseUint(cpuRevisionStr, 0, 64); err == nil {
					cpuInfo["stepping"] = fmt.Sprintf("r%dp%d", cpuVariant, cpuRevision)
				}
			}
		}
	}

	// Iterate over each processor and fetch additional information from /sys/devices/system/cpu
	cores := map[uint64]struct{}{}
	packages := map[uint64]struct{}{}
	cacheSizes := map[uint64]uint64{}
	for _, stanza := range procCpu {
		procID, err := strconv.ParseUint(stanza["processor"], 0, 64)
		if err != nil {
			continue
		}

		if coreID, ok := sysCpuInt(fmt.Sprintf("cpu%d/topology/core_id", procID)); ok {
			cores[coreID] = struct{}{}
		}

		if pkgID, ok := sysCpuInt(fmt.Sprintf("cpu%d/topology/physical_package_id", procID)); ok {
			packages[pkgID] = struct{}{}
		}

		// iterate over each cache this CPU can use
		i := 0
		for {
			if sharedList, ok := sysCpuList(fmt.Sprintf("cpu%d/cache/index%d/shared_cpu_list", procID, i)); ok {
				// we are scanning CPUs in order, so only count this cache if it's not shared with a
				// CPU that has already been scanned
				shared := false
				for sharedProcID := range sharedList {
					if sharedProcID < procID {
						shared = true
						break
					}
				}

				if !shared {
					if level, ok := sysCpuInt(fmt.Sprintf("cpu%d/cache/index%d/level", procID, i)); ok {
						if size, ok := sysCpuSize(fmt.Sprintf("cpu%d/cache/index%d/size", procID, i)); ok {
							cacheSizes[level] += size
						}
					}
				}
			} else {
				break
			}
			i++
		}
	}
	cpuInfo["cpu_pkgs"] = strconv.Itoa(len(packages))
	cpuInfo["cpu_cores"] = strconv.Itoa(len(cores))
	cpuInfo["cpu_logical_processors"] = strconv.Itoa(len(procCpu))
	cpuInfo["cache_size_l1"] = strconv.FormatUint(cacheSizes[1], 10)
	cpuInfo["cache_size_l2"] = strconv.FormatUint(cacheSizes[2], 10)
	cpuInfo["cache_size_l3"] = strconv.FormatUint(cacheSizes[3], 10)

	// cache_size uses the format '9216 KB'
	cpuInfo["cache_size"] = fmt.Sprintf("%d KB", (cacheSizes[1]+cacheSizes[2]+cacheSizes[3])/1024)

	// Count the number of NUMA nodes in /sys/devices/system/node
	nodes := 0
	if dirents, err := os.ReadDir("/sys/devices/system/node"); err == nil {
		for _, dirent := range dirents {
			if dirent.IsDir() && nodeNRegex.MatchString(dirent.Name()) {
				nodes++
			}
		}
	}
	cpuInfo["cpu_numa_nodes"] = strconv.Itoa(nodes)

	// ARM does not make the clock speed available
	// cpuInfo["mhz"]

	return cpuInfo, nil
}
