// This file is licensed under the MIT License.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright © 2015 Kentaro Kuribayashi <kentarok@gmail.com>
// Copyright 2014-present Datadog, Inc.

package platform

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"syscall"
	"unicode/utf16"
	"unsafe"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

// OSVERSIONINFOEXW contains operating system version information.
// From winnt.h (see https://docs.microsoft.com/en-us/windows/win32/api/winnt/ns-winnt-osversioninfoexw)
// This is used by https://docs.microsoft.com/en-us/windows/win32/devnotes/rtlgetversion
//
//nolint:revive
type OSVERSIONINFOEXW struct {
	dwOSVersionInfoSize uint32
	dwMajorVersion      uint32
	dwMinorVersion      uint32
	dwBuildNumber       uint32
	dwPlatformId        uint32
	szCSDVersion        [128]uint16
	wServicePackMajor   uint16
	wServicePackMinor   uint16
	wSuiteMask          uint16
	wProductType        uint8
	wReserved           uint8
}

var (
	modNetapi32          = windows.NewLazyDLL("Netapi32.dll")
	procNetServerGetInfo = modNetapi32.NewProc("NetServerGetInfo")
	procNetAPIBufferFree = modNetapi32.NewProc("NetApiBufferFree")
	ntdll                = windows.NewLazyDLL("Ntdll.dll")
	procRtlGetVersion    = ntdll.NewProc("RtlGetVersion")
	winbrand             = windows.NewLazyDLL("winbrand.dll")

	// ERROR_SUCCESS is the error returned in case of success
	ERROR_SUCCESS syscall.Errno
)

// see https://learn.microsoft.com/en-us/windows/win32/api/lmserver/nf-lmserver-netserverenum
//
//nolint:revive
const (
	// SV_TYPE_WORKSTATION is for all workstations.
	SV_TYPE_WORKSTATION = uint32(0x00000001)
	// SV_TYPE_SERVER is for all computers that run the Server service.
	SV_TYPE_SERVER = uint32(0x00000002)
	// SV_TYPE_SQLSERVER is for any server that runs an instance of Microsoft SQL Server.
	SV_TYPE_SQLSERVER = uint32(0x00000004)
	// SV_TYPE_DOMAIN_CTRL is for a server that is primary domain controller.
	SV_TYPE_DOMAIN_CTRL = uint32(0x00000008)
	// SV_TYPE_DOMAIN_BAKCTRL is for any server that is a backup domain controller.
	SV_TYPE_DOMAIN_BAKCTRL = uint32(0x00000010)
	// SV_TYPE_TIME_SOURCE is for any server that runs the Timesource service.
	SV_TYPE_TIME_SOURCE = uint32(0x00000020)
	// SV_TYPE_AFP is for any server that runs the Apple Filing Protocol (AFP) file service.
	SV_TYPE_AFP = uint32(0x00000040)
	// SV_TYPE_NOVELL is for any server that is a Novell server.
	SV_TYPE_NOVELL = uint32(0x00000080)
	// SV_TYPE_DOMAIN_MEMBER is for any computer that is LAN Manager 2.x domain member.
	SV_TYPE_DOMAIN_MEMBER = uint32(0x00000100)
	// SV_TYPE_PRINTQ_SERVER is for any computer that shares a print queue.
	SV_TYPE_PRINTQ_SERVER = uint32(0x00000200)
	// SV_TYPE_DIALIN_SERVER is for any server that runs a dial-in service.
	SV_TYPE_DIALIN_SERVER = uint32(0x00000400)
	// SV_TYPE_XENIX_SERVER is for any server that is a Xenix server.
	SV_TYPE_XENIX_SERVER = uint32(0x00000800)
	// SV_TYPE_SERVER_UNIX is for any server that is a UNIX server. This is the same as the SV_TYPE_XENIX_SERVER.
	SV_TYPE_SERVER_UNIX = SV_TYPE_XENIX_SERVER
	// SV_TYPE_NT is for a workstation or server.
	SV_TYPE_NT = uint32(0x00001000)
	// SV_TYPE_WFW is for any computer that runs Windows for Workgroups.
	SV_TYPE_WFW = uint32(0x00002000)
	// SV_TYPE_SERVER_MFPN is for any server that runs the Microsoft File and Print for NetWare service.
	SV_TYPE_SERVER_MFPN = uint32(0x00004000)
	// SV_TYPE_SERVER_NT is for any server that is not a domain controller.
	SV_TYPE_SERVER_NT = uint32(0x00008000)
	// SV_TYPE_POTENTIAL_BROWSER is for any computer that can run the browser service.
	SV_TYPE_POTENTIAL_BROWSER = uint32(0x00010000)
	// SV_TYPE_BACKUP_BROWSER is for a computer that runs a browser service as backup.
	SV_TYPE_BACKUP_BROWSER = uint32(0x00020000)
	// SV_TYPE_MASTER_BROWSER is for a computer that runs the master browser service.
	SV_TYPE_MASTER_BROWSER = uint32(0x00040000)
	// SV_TYPE_DOMAIN_MASTER is for a computer that runs the domain master browser.
	SV_TYPE_DOMAIN_MASTER = uint32(0x00080000)
	// SV_TYPE_SERVER_OSF is for a computer that runs OSF/1.
	SV_TYPE_SERVER_OSF = uint32(0x00100000)
	// SV_TYPE_SERVER_VMS is for a computer that runs Open Virtual Memory System (VMS).
	SV_TYPE_SERVER_VMS = uint32(0x00200000)
	// SV_TYPE_WINDOWS is for a computer that runs Windows.
	SV_TYPE_WINDOWS = uint32(0x00400000) /* Windows95 and above */
	// SV_TYPE_DFS is for a computer that is the root of Distributed File System (DFS) tree.
	SV_TYPE_DFS = uint32(0x00800000)
	// SV_TYPE_CLUSTER_NT is for server clusters available in the domain.
	SV_TYPE_CLUSTER_NT = uint32(0x01000000)
	// SV_TYPE_TERMINALSERVER is for a server running the Terminal Server service.
	SV_TYPE_TERMINALSERVER = uint32(0x02000000)
	// SV_TYPE_CLUSTER_VS_NT is for cluster virtual servers available in the domain.
	SV_TYPE_CLUSTER_VS_NT = uint32(0x04000000)
	// SV_TYPE_DCE is for a computer that runs IBM Directory and Security Services (DSS) or equivalent.
	SV_TYPE_DCE = uint32(0x10000000)
	// SV_TYPE_ALTERNATE_XPORT is for a computer that over an alternate transport.
	SV_TYPE_ALTERNATE_XPORT = uint32(0x20000000)
	// SV_TYPE_LOCAL_LIST_ONLY is for any computer maintained in a list by the browser. See the following Remarks section.
	SV_TYPE_LOCAL_LIST_ONLY = uint32(0x40000000)
	// SV_TYPE_DOMAIN_ENUM is for the primary domain.
	SV_TYPE_DOMAIN_ENUM = uint32(0x80000000)
	// SV_TYPE_ALL is for all servers. This is a convenience that will return all possible servers
	SV_TYPE_ALL = uint32(0xFFFFFFFF) /* handy for NetServerEnum2 */
)
const registryHive = "SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion"
const productNameKey = "ProductName"
const buildNumberKey = "CurrentBuildNumber"
const majorKey = "CurrentMajorVersionNumber"
const minorKey = "CurrentMinorVersionNumber"

func netServerGetInfo() (si SERVER_INFO_101, err error) {
	var outdata *byte
	// do additional work so that we don't panic() when the library's
	// not there (like in a container)
	if err = modNetapi32.Load(); err != nil {
		return
	}
	if err = procNetServerGetInfo.Find(); err != nil {
		return
	}
	status, _, err := procNetServerGetInfo.Call(uintptr(0), uintptr(101), uintptr(unsafe.Pointer(&outdata)))
	if status != uintptr(0) {
		return
	}
	defer procNetAPIBufferFree.Call(uintptr(unsafe.Pointer(outdata)))
	return platGetServerInfo(outdata), nil
}

func fetchOsDescription() (string, error) {
	err := winbrand.Load()
	if err == nil {
		// From https://stackoverflow.com/a/69462683
		procBrandingFormatString := winbrand.NewProc("BrandingFormatString")
		if procBrandingFormatString.Find() == nil {
			// Encode the string "%WINDOWS_LONG%" to UTF-16 and append a null byte for the Windows API
			magicString := utf16.Encode([]rune("%WINDOWS_LONG%" + "\x00"))
			os, _, err := procBrandingFormatString.Call(uintptr(unsafe.Pointer(&magicString[0])))
			defer syscall.LocalFree(syscall.Handle(os))
			if err == ERROR_SUCCESS {
				return windows.UTF16PtrToString((*uint16)(unsafe.Pointer(os))), nil
			}
		}
	}

	k, err := registry.OpenKey(registry.LOCAL_MACHINE,
		registryHive,
		registry.QUERY_VALUE)
	defer k.Close()
	if err == nil {
		os, _, err := k.GetStringValue(productNameKey)
		if err == nil {
			return os, nil
		}
	}

	return "(undetermined windows version)", err
}

func fetchWindowsVersion() (major uint64, minor uint64, build uint64, err error) {
	var osversion OSVERSIONINFOEXW
	status, _, _ := procRtlGetVersion.Call(uintptr(unsafe.Pointer(&osversion)))
	if status == 0 {
		major = uint64(osversion.dwMajorVersion)
		minor = uint64(osversion.dwMinorVersion)
		build = uint64(osversion.dwBuildNumber)
	} else {
		var regkey registry.Key
		regkey, err = registry.OpenKey(registry.LOCAL_MACHINE,
			registryHive,
			registry.QUERY_VALUE)
		defer regkey.Close()
		if err != nil {
			major, _, err = regkey.GetIntegerValue(majorKey)
			if err != nil {
				return
			}

			minor, _, err = regkey.GetIntegerValue(minorKey)
			if err != nil {
				return
			}

			var regbuild string
			regbuild, _, err = regkey.GetStringValue(buildNumberKey)
			if err != nil {
				return
			}
			build, err = strconv.ParseUint(regbuild, 10, 0)
		}
	}
	return
}

// GetArchInfo returns basic host architecture information
func GetArchInfo() (systemInfo map[string]string, err error) {
	systemInfo = map[string]string{}

	systemInfo["hostname"], _ = os.Hostname()

	if runtime.GOARCH == "amd64" {
		systemInfo["machine"] = "x86_64"
	} else {
		systemInfo["machine"] = runtime.GOARCH
	}

	systemInfo["os"], err = fetchOsDescription()

	maj, min, bld, err := fetchWindowsVersion()
	verstring := fmt.Sprintf("%d.%d.%d", maj, min, bld)
	systemInfo["kernel_release"] = verstring

	systemInfo["kernel_name"] = "Windows"

	// do additional work so that we don't panic() when the library's
	// not there (like in a container)
	family := "Unknown"
	si, sierr := netServerGetInfo()
	if sierr == nil {
		if (si.sv101_type&SV_TYPE_WORKSTATION) == SV_TYPE_WORKSTATION ||
			(si.sv101_type&SV_TYPE_SERVER) == SV_TYPE_SERVER {
			if (si.sv101_type & SV_TYPE_WORKSTATION) == SV_TYPE_WORKSTATION {
				family = "Workstation"
			} else if (si.sv101_type & SV_TYPE_SERVER) == SV_TYPE_SERVER {
				family = "Server"
			}
			if (si.sv101_type & SV_TYPE_DOMAIN_MEMBER) == SV_TYPE_DOMAIN_MEMBER {
				family = "Domain Joined " + family
			} else {
				family = "Standalone " + family
			}
		} else if (si.sv101_type & SV_TYPE_DOMAIN_CTRL) == SV_TYPE_DOMAIN_CTRL {
			family = "Domain Controller"
		} else if (si.sv101_type & SV_TYPE_DOMAIN_BAKCTRL) == SV_TYPE_DOMAIN_BAKCTRL {
			family = "Backup Domain Controller"
		}
	}
	systemInfo["family"] = family

	return
}
