package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/mux"
)

func CreateNotebook(resp http.ResponseWriter, req *http.Request) {
	contentLength := req.ContentLength
	fmt.Printf("Content Length Received : %v\n", contentLength)
	request := &CreateNotebookRequest{}
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatalf("Unable to read message from request : %v", err)
	}
	proto.Unmarshal(data, request)
	name := request.GetName()
	result := &{Message: "Hello " + name}
	response, err := proto.Marshal(result)
	if err != nil {
		log.Fatalf("Unable to marshal response : %v", err)
	}
	resp.Write(response)

}

func main() {
	log.Infoln("Starting notebook server")
	r := mux.NewRouter()
	r.HandleFunc("/create/notebook", CreateNotebook).Methods("POST")
	// r.HandleFunc("/note", CreateNote).Methods("POST")
	// r.HandleFunc("/notebook", GetNotebook).Methods("GET")
	// r.HandleFunc("/note", GetNote).Methods("GET")
	// r.HandleFunc("/note", UpdateNote).Methods("UPDATE")
	// r.HandleFunc("/note", DeleteNote).Methods("DELETE")

	server := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:8080",
		WriteTimeout: 2 * time.Second,
		ReadTimeout:  2 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}
