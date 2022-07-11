package cpu

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSysCpuInt(t *testing.T) {
	testingPrefix = t.TempDir()
	defer func() { testingPrefix = "" }()
	os.MkdirAll(filepath.Join(testingPrefix, filepath.FromSlash("sys/devices/system/cpu")), 0o777)
	path := filepath.Join(testingPrefix, filepath.FromSlash("sys/devices/system/cpu/somefile"))

	t.Run("zero", func(t *testing.T) {
		ioutil.WriteFile(path, []byte("0\n"), 0o666)
		got, ok := sysCpuInt("somefile")
		require.True(t, ok)
		require.Equal(t, uint64(0), got)
	})

	t.Run("dec", func(t *testing.T) {
		ioutil.WriteFile(path, []byte("20\n"), 0o666)
		got, ok := sysCpuInt("somefile")
		require.True(t, ok)
		require.Equal(t, uint64(20), got)
	})

	t.Run("hex", func(t *testing.T) {
		ioutil.WriteFile(path, []byte("0x20\n"), 0o666)
		got, ok := sysCpuInt("somefile")
		require.True(t, ok)
		require.Equal(t, uint64(32), got)
	})

	t.Run("invalid", func(t *testing.T) {
		ioutil.WriteFile(path, []byte("eleventy"), 0o666)
		_, ok := sysCpuInt("somefile")
		require.False(t, ok)
	})

	t.Run("missing", func(t *testing.T) {
		_, ok := sysCpuInt("nonexistent")
		require.False(t, ok)
	})
}

func TestSysCpuSize(t *testing.T) {
	testingPrefix = t.TempDir()
	defer func() { testingPrefix = "" }()
	os.MkdirAll(filepath.Join(testingPrefix, filepath.FromSlash("sys/devices/system/cpu")), 0o777)
	path := filepath.Join(testingPrefix, filepath.FromSlash("sys/devices/system/cpu/somefile"))

	t.Run("zero", func(t *testing.T) {
		ioutil.WriteFile(path, []byte("0\n"), 0o666)
		got, ok := sysCpuSize("somefile")
		require.True(t, ok)
		require.Equal(t, uint64(0), got)
	})

	t.Run("no-suffix", func(t *testing.T) {
		ioutil.WriteFile(path, []byte("20\n"), 0o666)
		got, ok := sysCpuSize("somefile")
		require.True(t, ok)
		require.Equal(t, uint64(20), got)
	})

	t.Run("K", func(t *testing.T) {
		ioutil.WriteFile(path, []byte("20K\n"), 0o666)
		got, ok := sysCpuSize("somefile")
		require.True(t, ok)
		require.Equal(t, uint64(20*1024), got)
	})

	t.Run("M", func(t *testing.T) {
		ioutil.WriteFile(path, []byte("20M"), 0o666)
		got, ok := sysCpuSize("somefile")
		require.True(t, ok)
		require.Equal(t, uint64(20*1024*1024), got)
	})

	t.Run("G", func(t *testing.T) {
		ioutil.WriteFile(path, []byte("20G"), 0o666)
		got, ok := sysCpuSize("somefile")
		require.True(t, ok)
		require.Equal(t, uint64(20*1024*1024*1024), got)
	})

	t.Run("invalid", func(t *testing.T) {
		ioutil.WriteFile(path, []byte("eleventy"), 0o666)
		_, ok := sysCpuSize("somefile")
		require.False(t, ok)
	})

	t.Run("missing", func(t *testing.T) {
		_, ok := sysCpuSize("nonexistent")
		require.False(t, ok)
	})
}

func TestSysCpuList(t *testing.T) {
	testingPrefix = t.TempDir()
	defer func() { testingPrefix = "" }()
	os.MkdirAll(filepath.Join(testingPrefix, filepath.FromSlash("sys/devices/system/cpu")), 0o777)
	path := filepath.Join(testingPrefix, filepath.FromSlash("sys/devices/system/cpu/somefile"))

	t.Run("empty", func(t *testing.T) {
		ioutil.WriteFile(path, []byte("\n"), 0o666)
		got, ok := sysCpuList("somefile")
		require.True(t, ok)
		require.Equal(t, map[uint64]struct{}{}, got)
	})

	t.Run("single", func(t *testing.T) {
		ioutil.WriteFile(path, []byte("20\n"), 0o666)
		got, ok := sysCpuList("somefile")
		require.True(t, ok)
		require.Equal(t, map[uint64]struct{}{20: {}}, got)
	})

	t.Run("range", func(t *testing.T) {
		ioutil.WriteFile(path, []byte("5-8\n"), 0o666)
		got, ok := sysCpuList("somefile")
		require.True(t, ok)
		require.Equal(t, map[uint64]struct{}{
			5: {},
			6: {},
			7: {},
			8: {},
		}, got)
	})

	t.Run("combo", func(t *testing.T) {
		ioutil.WriteFile(path, []byte("1,5-8,10\n"), 0o666)
		got, ok := sysCpuList("somefile")
		require.True(t, ok)
		require.Equal(t, map[uint64]struct{}{
			1:  {},
			5:  {},
			6:  {},
			7:  {},
			8:  {},
			10: {},
		}, got)
	})

	t.Run("invalid", func(t *testing.T) {
		ioutil.WriteFile(path, []byte("eleventy"), 0o666)
		_, ok := sysCpuList("somefile")
		require.False(t, ok)
	})

	t.Run("missing", func(t *testing.T) {
		_, ok := sysCpuList("nonexistent")
		require.False(t, ok)
	})
}