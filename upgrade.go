package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

const (
	TmpPostfix = ".tmp"
)

func Upgrade(url string, location string, name string) {
	download(url, location, name)
	startCliInit(location, name)
	os.Exit(0)
}

func InitDownloaded() {
	renameCurrentBinary()
	os.Exit(0)
}

func download(url string, location string, name string) {
	if _, err := os.Stat(location + name + TmpPostfix); err == nil {
		err = os.Remove(location + name + TmpPostfix)
		if err != nil {
			log.Fatalf("Failed to remove previos %s binary : %s", TmpPostfix, err.Error())
		}
	}
	out, err := os.OpenFile(location+name+TmpPostfix, os.O_CREATE|os.O_WRONLY, 0700)
	if err != nil {
		log.Fatalf("Failed to write new binary on disk: %s", err.Error())
	}
	defer out.Close()
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Failed to download new binary: %s", err.Error())
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Fatalf("Failed to download new binary: %s", err.Error())
	}
}

func startCliInit(location string, name string) {
	if runtime.GOOS != "Windows" {
		if location == "" {
			location = "./"
		}
	}
	cmd := exec.Command(location+name+TmpPostfix, "-downloaded=true")
	err := cmd.Run()
	if err != nil {
		log.Println(err.Error())
	}
}

//remove postfix from the name of executing binary
func renameCurrentBinary() {
	executable, err := os.Executable()
	if err != nil {
		log.Fatalf("failed to get executable : %s \n", err.Error())
	}
	separator := "/"
	if runtime.GOOS == "windows" {
		separator = "\\"
	}
	execPath := strings.SplitAfter(executable, separator)
	execFile := execPath[len(execPath)-1]
	execDir := strings.Join(execPath[:len(execPath)-1], "")
	log.Println("Trying to rename myself")
	err = os.Rename(execDir+execFile, execDir+execFile[0:len(execFile)-len(TmpPostfix)])
	if err != nil {
		log.Fatalf("Error while renaming file: %s\n", err.Error())
	}
	log.Println("Renaming succeed")
}
