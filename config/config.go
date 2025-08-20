package config

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
)

var (
	// ToolDir is the template to find the root tool directory
	ToolDir = "{{RootDir}}/.tool"
	// RootDir is the template to find the root directory when executing
	RootDir = "{{GitRoot}}"
	// TmpDir is the template to find the an alternate TempDir, if empty defaults to system temp dir
	TmpDir = ""

	// OS is the OS name to request for commands that require OS name
	OS = runtime.GOOS
	// Arch is the architecture to request for commands that require architecture
	Arch = runtime.GOARCH

	// Debug whether to output debug logging and perform additional diagnostic work
	Debug = false
	// Trace enables Debug and enables even more verbose logging
	Trace = false

	// CI indicates running in a CI environment
	CI = false
	// Cleanup whether to remove temporary files and downloads
	Cleanup = true
)

func init() {
	_, _ = fmt.Fprintf(os.Stderr, "initializing config with env: %v\n", os.Environ())
	Trace, _ = strconv.ParseBool(Env("TRACE", "false"))
	Debug, _ = strconv.ParseBool(Env("DEBUG",
		Env("ACTIONS_RUNNER_DEBUG", strconv.FormatBool(Trace))))
	CI, _ = strconv.ParseBool(Env("CI", "false"))
	Cleanup = !Debug && !CI
}
