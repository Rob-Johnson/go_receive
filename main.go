package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/Rob-Johnson/goreceive/deploy"
	"io/ioutil"
	"net/http"
)

//only need the after field from Github's post receive hook
type Payload struct {
	After string
}

func sendResponse(w http.ResponseWriter, status int, data string) {
	w.WriteHeader(status)
	w.Write([]byte(fmt.Sprintf(data)))
}

func ReceiveHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		sendResponse(w, 405, fmt.Sprintf("Need a POST, not a %s.", r.Method))
	}

	if r.Body == nil {
		sendResponse(w, 400, "No Data Received")
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		sendResponse(w, 400, err.Error())
		return
	}

	var load Payload
	err = json.Unmarshal(body, &load)
	if err != nil {
		sendResponse(w, 400, err.Error())
		return
	}
	if load.After == "" {
		sendResponse(w, 400, "Empty string")
		return
	}

	err = goreceive.CheckEnv()
	if err != nil {
		sendResponse(w, 500, fmt.Sprintf("%s", "Problem updating codebase"))
	}

	err = goreceive.RedeployCodebase(load.After)
	if err != nil {
		sendResponse(w, 500, fmt.Sprintf("%s", "Problem updating codebase"))
	}
}

func main() {
	var port = flag.Int("port", 8080, "Port to listen on")
	flag.Parse()
	http.HandleFunc("/", ReceiveHandler)
	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}
