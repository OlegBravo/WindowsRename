package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

type config struct {
	IsDownloaded bool
	LogFile      string
}

func main() {

	c := readArguments(os.Args)

	//cli -> download cli.tmp
	//cli -> run cli.tmp /flag
	//cli -> exit
	//cli.tmp -> myCopy cli.tmp cli

	if c.IsDownloaded {
		InitForDownloaded(c)
	} else {
		DownloadNewVersion(c)
	}

}

func InitForDownloaded(c config) {
	f, err := os.OpenFile(c.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	executable, err := os.Executable()
	if err != nil {
		log.Fatalf("failed to get executable : %s \n", err.Error())
	}
	execPath := strings.SplitAfter(executable, "\\")
	execFile := execPath[len(execPath)-1]
	execDir := strings.Join(execPath[:len(execPath)-1], "\\")
	log.Println("Trying to rename myself")
	err = os.Rename(execDir+execFile, execDir+execFile[0:len(execFile)-4])
	if err != nil {
		log.Fatalf("Error while renaming file: %s\n", err.Error())
	}
	log.Println("Renaming succeed")

	//rename myself back to normal value
}

func DownloadNewVersion(c config) {
	executable, _ := os.Executable()
	execPath := strings.SplitAfter(executable, "\\")
	execFile := execPath[len(execPath)-1]
	execDir := strings.Join(execPath[:len(execPath)-1], "\\")

	log.Println("Pretending download is happening")
	myCopy(execDir+execFile, execFile+".tmp")

	cmd := exec.Command(execFile+".tmp", "-downloaded=true", "logfile="+c.LogFile)
	err := cmd.Run()
	if err != nil {
		log.Println(err.Error())
	}

	log.Println("Exiting current ")
	os.Exit(0)

}

func myCopy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func readArguments(a []string) config {
	isDownloaded := flag.Bool("downloaded", false, "")
	logFile := flag.String("logfile", "C:\\Users\\79169\\go\\src\\WindowsRename\\logfile.log", "")
	flag.Parse()
	return config{*isDownloaded, *logFile}
}
