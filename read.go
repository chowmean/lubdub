package main

import (
    //"fmt"
    "strconv"
    "time"
    "io/ioutil"
    "os"
    "path/filepath"
    "regexp"
    "./httpClass"
    "github.com/golang/glog"
)

func check(e error) {
    if e != nil {
	glog.Info("Error reading file" + string(e.Error()))
    }
    glog.Flush()
}

func formatSend(data string){
	argsWithoutProg := os.Args[1:]
	url :=  argsWithoutProg[0]
	hostname, err := os.Hostname()
	check(err)
	content := httpClass.Content{
		Content: data,
		ID: 1,
		Hostname:hostname,
	}
	client := httpClass.BasicAuthClient("Token")
	client.PostStatus(&content,url)
}

func readCPU(){
	cpustat,err := ioutil.ReadFile("/proc/stat")
        check(err)
	go formatSend(string(cpustat))
}

func readPROC(file string){
	procdata,err := ioutil.ReadFile(file)
	check(err)
	go formatSend(string(procdata))
}


func main(){
	argsWithoutProg := os.Args[1:]
	ttl :=  argsWithoutProg[1]
	for{
		go readCPU()
		searchDir := "/proc/"
		filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
			match, _ := regexp.MatchString("/proc/([0-9]+)/status", path)
			if(match){
				go readPROC(path)
			}
	        	return nil
		})
	i, err := strconv.Atoi(ttl)
	check(err)
	time.Sleep(time.Duration(i) * time.Second)
	}
}
