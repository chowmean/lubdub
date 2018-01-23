package main

import (
    "fmt"
    "io/ioutil"
    "os"
    "path/filepath"
    "regexp"
    "./httpClass"
)

func check(e error) {
    if e != nil {
        fmt.Print("error")
    }
}

func format_send(data string){
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
	cpu_stat,err := ioutil.ReadFile("/proc/stat")
        check(err)
	go format_send(string(cpu_stat))
}

func readPROC(file string){
	proc_data,err := ioutil.ReadFile(file)
	check(err)
	go format_send(string(proc_data))
}


func main(){
	go readCPU()
	searchDir := "/proc/"
	filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		match, _ := regexp.MatchString("/proc/([0-9]+)/status", path)
		if(match){
			go readPROC(path)
		}
	        return nil
	})
}
