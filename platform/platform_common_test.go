package platform

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/DataDog/gohai/utils"
)

func TestGetPythonVersion(t *testing.T) {
	t.Run("valid Python version", func(t *testing.T) {
		pythonV, err := getPythonVersion(pythonVersionFakeExecCmd("TestGetPythonVersionCmd", "valid"))
		assert.Nil(t, err)
		assert.Equal(t, "3.8.9", pythonV)
	})

	t.Run("valid Python version (windows)", func(t *testing.T) {
		pythonV, err := getPythonVersion(pythonVersionFakeExecCmd("TestGetPythonVersionCmd", "valid-windows"))
		assert.Nil(t, err)
		assert.Equal(t, "3.8.9", pythonV)
	})

	t.Run("Python not present", func(t *testing.T) {
		pythonV, err := getPythonVersion(pythonVersionFakeExecCmd("TestGetPythonVersionCmd", "not-present"))
		assert.NotNil(t, err)
		assert.Equal(t, "", pythonV)
	})

	t.Run("invalid Python version", func(t *testing.T) {
		pythonV, err := getPythonVersion(pythonVersionFakeExecCmd("TestGetPythonVersionCmd", "invalid"))
		assert.NotNil(t, err)
		assert.Equal(t, "", pythonV)
	})
}

// TestGetHostnameShellCmd is a method that is called as a substitute for a shell command,
// the GO_TEST_PROCESS flag ensures that if it is called as part of the test suite, it is skipped.
func TestGetPythonVersionCmd(t *testing.T) {
	if os.Getenv("GO_TEST_PROCESS") != "1" {
		return
	}

	testRunName, cmdList := utils.ParseFakeExecCmdArgs()

	assert.EqualValues(t, []string{"python", "-V"}, cmdList)

	switch testRunName {
	case "valid":
		fmt.Fprintf(os.Stdout, "Python 3.8.9\n")
		os.Exit(0)
	case "valid-windows":
		fmt.Fprintf(os.Stdout, "Python 3.8.9\r\n")
		os.Exit(0)
	case "not-present":
		fmt.Fprintf(os.Stdout, "")
		fmt.Fprintf(os.Stderr, "command not found: python")
		os.Exit(127)
	case "invalid":
		fmt.Fprintf(os.Stdout, "giberrish")
		os.Exit(0)
	}
}

func pythonVersionFakeExecCmd(testName string, testRunName string) utils.ExecCmdFunc {
	return func(command string, args ...string) *exec.Cmd {
		return utils.FakeExecCmd(testName, testRunName, command, args...)
	}
}
