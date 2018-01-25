package main

import (
	//"fmt"
	"./httpClass"
	"github.com/golang/glog"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"
	"math/rand"
)

func check(e error) {
	if e != nil {
		glog.Info("Error reading file" + string(e.Error()))
	}
	glog.Flush()
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
    b := make([]byte, n)
    for i := range b {
        b[i] = letterBytes[rand.Intn(len(letterBytes))]
    }
    return string(b)
}

func formatSend(data string, id string) {
	argsWithoutProg := os.Args[1:]
	url := argsWithoutProg[0]
	hostname, err := os.Hostname()
	check(err)
	content := httpClass.Content{
		Content:  data,
		ID:       id,
		Hostname: hostname,
	}
	client := httpClass.BasicAuthClient("Token")
	client.PostStatus(&content, url)
}

func readCPU(id string) {
	cpustat, err := ioutil.ReadFile("/proc/stat")
	check(err)
	go formatSend(string(cpustat), id)
}

func readPROC(file string, id string) {
	procdata, err := ioutil.ReadFile(file)
	check(err)
	go formatSend(string(procdata), id)
}

func main() {
	argsWithoutProg := os.Args[1:]
	ttl := argsWithoutProg[1]
	for {
		id := RandStringBytes(18)
		go readCPU(id)
		searchDir := "/proc/"
		filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
			match, _ := regexp.MatchString("/proc/([0-9]+)/status", path)
			if match {
				go readPROC(path,id)
			}
			return nil
		})
		i, err := strconv.Atoi(ttl)
		check(err)
		time.Sleep(time.Duration(i) * time.Second)
	}
}
