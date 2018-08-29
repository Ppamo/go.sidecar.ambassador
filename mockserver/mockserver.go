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
	fmt.Printf("+ RH: %s::%s\n", r.Method, r.RequestURI)
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

func getResponse(serviceName string) (string, error) {
	response, err := loadServiceMock(serviceName)
	if err != nil {
		return get404ErrorResponse(), err
	}
	mock := new(Mock)
	err = json.Unmarshal([]byte(response), &mock)
	if err != nil {
		return "{}", fmt.Errorf("JSon UnMarshal failed\n%s", err)
	}
	body, err := json.Marshal(mock.Body)
	if err != nil {
		return "{}", err
	}
	return string(body), nil
}

func getRequest(w http.ResponseWriter, r *http.Request) (*Request, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil { return nil, err }
	if len(body) == 0 { return nil, errors.New("Body is empty!") }
	defer r.Body.Close()
	request := new(Request)
	err = json.Unmarshal(body, &request)
	if err != nil { return nil, fmt.Errorf("JSon UnMarshal failed\n%s", err) }
	return request, nil
}

func requestBroker(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("+ RB: %s::%s\n", r.Method, r.RequestURI)
	request,err := getRequest(w, r)
	if (err != nil) {
		fmt.Printf("- ERROR: %v\n", err)
		fmt.Fprintf(w, getErrorResponse())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	response, err := getResponse(request.Data["service_name"].(string))
	fmt.Fprintf(w, response)
}
