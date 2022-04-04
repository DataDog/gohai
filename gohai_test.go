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
	} `json:"cpu"`
	Filesystem []struct {
		KbSize    string `json:"kb_size"`
		MountedOn string `json:"mounted_on"`
		Name      string `json:"name"`
	} `json:"filesystem"`
	Memory struct {
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
		KernelVersion string `json:"kernel_version"`
		Machine       string `json:"machine"`
		Os            string `json:"os"`
		Processor     string `json:"processor"`
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

	if assert.NotEmpty(t, payload.Filesystem) {
		if runtime.GOOS != "windows" {
			assert.NotEmpty(t, payload.Filesystem[0].KbSize, 0)
		}
		assert.NotEmpty(t, payload.Filesystem[0].MountedOn, 0)
		assert.NotEmpty(t, payload.Filesystem[0].Name, 0)
	}
	if runtime.GOOS != "windows" {
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
	if runtime.GOOS != "windows" {
		assert.NotEmpty(t, payload.Platform.KernelVersion)
	}
	assert.NotEmpty(t, payload.Platform.Machine)
	assert.NotEmpty(t, payload.Platform.Os)
	if runtime.GOOS != "windows" {
		assert.NotEmpty(t, payload.Platform.Processor)
	}
}
