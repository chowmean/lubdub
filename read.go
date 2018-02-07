package main

import (
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


func randomInt(min, max int) int {
    return min + rand.Intn(max-min)
}

func check(e error) {
	if e != nil {
		glog.Info("Error reading file" + string(e.Error()))
	}
	glog.Flush()
}

func randomString(len int) string {
    bytes := make([]byte, len)
    for i := 0; i < len; i++ {
        bytes[i] = byte(randomInt(65, 90))
    }
    return string(bytes)
}

func formatSend(data string, id string, typeS string) {
	argsWithoutProg := os.Args[1:]
	url := argsWithoutProg[0]
	token := argsWithoutProg[1]
	hostname, err := os.Hostname()
	check(err)
	content := httpClass.Content{
		Content:  data,
		ID:       id,
		Hostname: hostname,
		ApiAccessToken: token,
		Type: typeS,
	}
	client := httpClass.BasicAuthClient("Token")
	client.PostStatus(&content, url)
}

func readCPU(id string) {
	cpustat, err := ioutil.ReadFile("/proc/stat")
	check(err)
	go formatSend(string(cpustat), id, "CPU")
}

func readPROC(file string, id string) {
	procdata, err := ioutil.ReadFile(file)
	check(err)
	go formatSend(string(procdata), id, "CPU PROCESS")
}

func readMemInfo(id string){
	meminfo, err := ioutil.ReadFile("/proc/meminfo")
	check(err)
	go formatSend(string(meminfo), id, "MEMORY")
}

func readCPUInfo(id string){
	cpuinfo,err := ioutil.ReadFile("/proc/cpuinfo")
	check(err)
	go formatSend(string(cpuinfo), id, "CPUINFO")
}

func main() {
	rand.Seed(time.Now().UnixNano())
	argsWithoutProg := os.Args[1:]
	ttl := argsWithoutProg[2]
	process_info := argsWithoutProg[3]
	for {
		id := randomString(18)
		go readCPU(id)
		go readMemInfo(id)
		go readCPUInfo(id)
		searchDir := "/proc/"
		if process_info == "1"{
			filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
				match, _ := regexp.MatchString("/proc/([0-9]+)/status", path)
				if match {
					go readPROC(path,id)
				}
				return nil
			})
		}
		i, err := strconv.Atoi(ttl)
		check(err)
		time.Sleep(time.Duration(i) * time.Second)
	}
}
