package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

//only need the after field from Github's post receive hook
type Payload struct {
	after string
}

func sendResponse(w http.ResponseWriter, status int, data string) {
	w.WriteHeader(status)
	w.Write([]byte(fmt.Sprintf(data)))
}

func ReceiveHandler(w http.ResponseWriter, r *http.Request) {

	//not interested if you're not POST'ing
	//return 405 Method Not Allowed
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

	//try and convert the json to a Payload struct
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

	//TODO: redeploy the codebase here
}

func main() {
	http.HandleFunc("/", ReceiveHandler)
	http.ListenAndServe(":8080", nil)
}
