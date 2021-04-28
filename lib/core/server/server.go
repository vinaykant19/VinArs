package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	cnf "../../../configuration"
	router "../router"
)
/*
 * Bootstrap functions using by base framework
 * @author: Vinaykant (vinaykantsahu@gmail.com)
 */
func BootstrapService(env string) {
	configuration := cnf.Configuration{}
	err := cnf.SetConfig(&configuration, env)

	allCalls := make(map[string]map[string][]string)
	readCall := make(map[string][]string)
	writeCall := make(map[string][]string)
	if err == nil {
		//get Looger
		setLogger(configuration.Logpath)
		log.Println("Reading routes")
		allCalls["read"] =router.CallParser("read", readCall)
		allCalls["write"] =router.CallParser("write", writeCall)
		//fmt.Println(allCalls)
		port := configuration.Port
		http.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
			router.Call(response, request, allCalls, &configuration)
		})
		err = http.ListenAndServe(":"+port, nil)
	}
	if err != nil {
		fmt.Println(err.Error())
	}
}

//
// New returns a new logger instance.
//
func setLogger(logFilePath string){
	date := time.Now().Format("2006-01-02")
	fileName := logFilePath + date + ".log"

	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	log.SetOutput(f)
}