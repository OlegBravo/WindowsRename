package main

import (
	"flag"
	"os"
	"runtime"
	"strings"
)

type config struct {
	IsDownloaded bool
}

const (
	DEFAULT_CLI_URL_WIN = "https://github.com/OlegBravo/WindowsRename/raw/master/WindowsRename.exe"
	DEFAULT_CLI_BIN_WIN = "WindowsRename.exe"

	DEFAULT_CLI_URL_DARWIN = "https://github.com/OlegBravo/WindowsRename/raw/master/WindowsRename.Darwin"
	DEFAULT_CLI_BIN_DARWIN = "WindowsRename.Darwin"

	DEFAULT_CLI_URL_UBUNTU = "https://github.com/OlegBravo/WindowsRename/raw/master/WindowsRename.Ubuntu"
	DEFAULT_CLI_BIN_UBUNTU = "WindowsRename.Ubuntu"
)

func main() {

	c := readArguments(os.Args)

	executable, _ := os.Executable()
	execPath := strings.SplitAfter(executable, "\\")
	execDir := strings.Join(execPath[:len(execPath)-1], "\\")
	println(execDir)

	if c.IsDownloaded {
		InitDownloaded()
	} else {
		if runtime.GOOS == "windows" {
			Upgrade(DEFAULT_CLI_URL_WIN,
				execDir, DEFAULT_CLI_BIN_WIN)
		} else if runtime.GOOS == "darwin" {
			Upgrade(DEFAULT_CLI_URL_DARWIN,
				"", DEFAULT_CLI_BIN_DARWIN)
		} else if runtime.GOOS == "ubuntu" {
			Upgrade(DEFAULT_CLI_URL_UBUNTU,
				"", DEFAULT_CLI_BIN_UBUNTU)
		}

	}

}

func readArguments(a []string) config {
	isDownloaded := flag.Bool("downloaded", false, "")
	flag.Parse()
	return config{*isDownloaded}
}
