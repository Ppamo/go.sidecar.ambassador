package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

var server *http.Server

func main() {
	fmt.Printf("=> starting mockserver process at port 8080\n")
	server = &http.Server{Addr: "0.0.0.0:8080"}
	http.HandleFunc("/", requestHandler)
	http.HandleFunc("/requestbroker", requestBroker)
	e := server.ListenAndServe()
	if e != nil {
		fmt.Printf("Error starting the server\n%v\n", e)
	}
}


func requestHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "{\"status\": \"ok\"}")
}

func getErrorResponse() string {
	return `{
	"StatusCode": 500,
	"StatusDescription": "internal error"
}`
}

func get404ErrorResponse() string {
	return `{
	"StatusCode": 404,
	"StatusDescription": "not found"
}`
}

func getMockPath(serviceName string) string {
	return fmt.Sprintf("/mocks/%s.mock.json", serviceName)
}

func loadServiceMock(serviceName string) (string, error) {
	fmt.Printf("=> Path: %s\n", getMockPath(serviceName))
	file, err := os.Open(getMockPath(serviceName))
	if err != nil {return "", err }
	defer file.Close()
	buffer := new(bytes.Buffer)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Fprintf(buffer, "%s", scanner.Text())
	}
	if err := scanner.Err(); err != nil { return "", err }
	return buffer.String(), nil
}

func getResponse(serviceName string) string {
	response, err := loadServiceMock(serviceName)
	if err != nil {
		return get404ErrorResponse()
	}
	return response
}

func getRequest(w http.ResponseWriter, r *http.Request) (*Request, error) {
	if r.Body == nil { return nil, errors.New("La peticion es vacia") }
	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)
	request := new(Request)
	err := json.Unmarshal(body, &request)
	if err != nil { return nil, fmt.Errorf("JSon UnMarshal failed\n%s", err) }
	return request, nil
}

func requestBroker(w http.ResponseWriter, r *http.Request) {
	request,err := getRequest(w, r)
	if (err != nil) {
		fmt.Printf("=> %v\n", request)
		fmt.Fprintf(w, getErrorResponse())
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, getResponse(request.Data["service_name"].(string)))
}
