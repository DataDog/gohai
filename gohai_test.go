package main

import (
	"encoding/json"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSelectedCollectors_String(t *testing.T) {
	sc := &SelectedCollectors{
		"foo": struct{}{},
		"bar": struct{}{},
	}
	assert.Equal(t, "[bar foo]", sc.String())
}

// gohaiPayload defines the format we expect the gohai information
// to be in.
type gohaiPayload struct {
	CPU struct {
		CPUCores             string `json:"cpu_cores"`
		CPULogicalProcessors string `json:"cpu_logical_processors"`
		Family               string `json:"family"`
		Mhz                  string `json:"mhz"`
		Model                string `json:"model"`
		ModelName            string `json:"model_name"`
		Stepping             string `json:"stepping"`
		VendorID             string `json:"vendor_id"`

		// On Windows, we report additional fields
		CacheSizeL1  string `json:"cache_size_l1"`
		CacheSizeL2  string `json:"cache_size_l2"`
		CacheSizeL3  string `json:"cache_size_l3"`
		CPUNumaNodes string `json:"cpu_numa_nodes"`
		CPUPkgs      string `json:"cpu_pkgs"`
	} `json:"cpu"`
	Filesystem []struct {
		KbSize string `json:"kb_size"`
		// MountedOn can be empty on Windows
		MountedOn string `json:"mounted_on"`
		Name      string `json:"name"`
	} `json:"filesystem"`
	Memory struct {
		// SwapTotal is not reported on Windows
		SwapTotal string `json:"swap_total"`
		Total     string `json:"total"`
	} `json:"memory"`
	Network struct {
		Interfaces []struct {
			Ipv4        []string `json:"ipv4"`
			Ipv6        []string `json:"ipv6"`
			Ipv6Network string   `json:"ipv6-network"`
			Macaddress  string   `json:"macaddress"`
			Name        string   `json:"name"`
			Ipv4Network string   `json:"ipv4-network"`
		} `json:"interfaces"`
		Ipaddress   string `json:"ipaddress"`
		Ipaddressv6 string `json:"ipaddressv6"`
		Macaddress  string `json:"macaddress"`
	} `json:"network"`
	Platform struct {
		Gooarch       string `json:"GOOARCH"`
		Goos          string `json:"GOOS"`
		GoV           string `json:"goV"`
		Hostname      string `json:"hostname"`
		KernelName    string `json:"kernel_name"`
		KernelRelease string `json:"kernel_release"`
		// KernelVersion is not reported on Windows
		KernelVersion string `json:"kernel_version"`
		Machine       string `json:"machine"`
		Os            string `json:"os"`
		Processor     string `json:"processor"`
		// On Windows, we report additional fields
		Family string `json:"family"`
	} `json:"platform"`
}

func TestGohaiSerialization(t *testing.T) {
	gohai, err := Collect()

	assert.NoError(t, err)

	gohaiJson, err := json.Marshal(gohai)
	assert.NoError(t, err)

	var payload gohaiPayload
	assert.NoError(t, json.Unmarshal(gohaiJson, &payload))

	assert.NotEmpty(t, payload.CPU.CPUCores)
	assert.NotEmpty(t, payload.CPU.CPULogicalProcessors)
	assert.NotEmpty(t, payload.CPU.Family)
	assert.NotEmpty(t, payload.CPU.Mhz)
	assert.NotEmpty(t, payload.CPU.Model)
	assert.NotEmpty(t, payload.CPU.ModelName)
	assert.NotEmpty(t, payload.CPU.Stepping)
	assert.NotEmpty(t, payload.CPU.VendorID)

	if runtime.GOOS == "windows" {
		// Additional fields that we report on Windows
		assert.NotEmpty(t, payload.CPU.CacheSizeL1)
		assert.NotEmpty(t, payload.CPU.CacheSizeL2)
		assert.NotEmpty(t, payload.CPU.CacheSizeL3)
		assert.NotEmpty(t, payload.CPU.CPUNumaNodes)
		assert.NotEmpty(t, payload.CPU.CPUPkgs)
	}

	if assert.NotEmpty(t, payload.Filesystem) {
		if runtime.GOOS != "windows" {
			// On Windows, MountedOn can be empty
			assert.NotEmpty(t, payload.Filesystem[0].MountedOn, 0)
		}
		assert.NotEmpty(t, payload.Filesystem[0].KbSize, 0)
		assert.NotEmpty(t, payload.Filesystem[0].Name, 0)
	}
	if runtime.GOOS != "windows" {
		// Not reported on Windows
		assert.NotEmpty(t, payload.Memory.SwapTotal)
	}
	assert.NotEmpty(t, payload.Memory.Total)

	if assert.NotEmpty(t, payload.Network.Interfaces) {
		assert.NotEmpty(t, payload.Network.Interfaces[0].Name)
		assert.NotEmpty(t, payload.Network.Interfaces[0].Macaddress)
		if len(payload.Network.Interfaces[0].Ipv4) == 0 {
			assert.NotEmpty(t, payload.Network.Interfaces[0].Ipv6)
			assert.NotEmpty(t, payload.Network.Interfaces[0].Ipv6Network)
		} else {
			assert.NotEmpty(t, payload.Network.Interfaces[0].Ipv4)
			assert.NotEmpty(t, payload.Network.Interfaces[0].Ipv4Network)
		}
	}
	assert.NotEmpty(t, payload.Network.Ipaddress)
	assert.NotEmpty(t, payload.Network.Ipaddressv6)
	assert.NotEmpty(t, payload.Network.Macaddress)

	assert.NotEmpty(t, payload.Platform.Gooarch)
	assert.NotEmpty(t, payload.Platform.Goos)
	assert.NotEmpty(t, payload.Platform.GoV)
	assert.NotEmpty(t, payload.Platform.Hostname)
	assert.NotEmpty(t, payload.Platform.KernelName)
	assert.NotEmpty(t, payload.Platform.KernelRelease)
	assert.NotEmpty(t, payload.Platform.Machine)
	assert.NotEmpty(t, payload.Platform.Os)
	if runtime.GOOS != "windows" {
		// Not reported on Windows
		assert.NotEmpty(t, payload.Platform.KernelVersion)
		assert.NotEmpty(t, payload.Platform.Processor)
	} else {
		// Additional fields that we report on Windows
		assert.NotEmpty(t, payload.Platform.Family)
	}
}
