package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

const (
	TmpPostfix = ".tmp"
)

func Upgrade(url string, location string, name string) {
	download(url, location, name)
	runNew(location, name)
	os.Exit(0)
}

func InitDownloaded() {
	renameMyBinary()
	os.Exit(0)
}

func download(url string, location string, name string) {
	out, err := os.Create(location + name + TmpPostfix)
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

func runNew(location string, name string) {
	cmd := exec.Command(location+name+TmpPostfix, "-downloaded=true")
	err := cmd.Run()
	if err != nil {
		log.Println(err.Error())
	}
}

//remove postfix from the name of executing binary
func renameMyBinary() {
	executable, err := os.Executable()
	if err != nil {
		log.Fatalf("failed to get executable : %s \n", err.Error())
	}
	execPath := strings.SplitAfter(executable, "\\")
	execFile := execPath[len(execPath)-1]
	execDir := strings.Join(execPath[:len(execPath)-1], "\\")
	log.Println("Trying to rename myself")
	err = os.Rename(execDir+execFile, execDir+execFile[0:len(execFile)-len(TmpPostfix)])
	if err != nil {
		log.Fatalf("Error while renaming file: %s\n", err.Error())
	}
	log.Println("Renaming succeed")
}
