package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/golang/protobuf/ptypes"

	// "time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/google/uuid"
)

const requiredProperty = "%s cannot be empty"

// NotebookRepo is an in memory db holding all notebooks
type NotebookRepo struct {
	notebooks map[string]Notebook
}

// NewNotebookRepo returns a reference to a NotebookRepo object
func NewNotebookRepo() *NotebookRepo {
	return &NotebookRepo{
		notebooks: make(map[string]Notebook),
	}
}

type tags map[string][]string

type Notebook struct {
	notes map[string]*Note
	// tags is a map with names as the key with values holding
	// a slice of string note IDs
	tags tags
}

// CreateNotebook takes a response body attempting to deserialise it to a
// CreateNotebookRequest object
func (n *NotebookRepo) CreateNotebook(w http.ResponseWriter, r *http.Request) {
	body := &CreateNotebookRequest{}
	if err := jsonpb.Unmarshal(r.Body, body); err != nil {
		log.Fatalf("Unable to unmarshal message from request: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if _, ok := n.notebooks[body.GetName()]; ok {
		http.Error(w, fmt.Sprintf("Notebook with name '%s' already exists", body.GetName()), http.StatusConflict)
		return
	}
	n.notebooks[body.GetName()] = Notebook{
		notes: make(map[string]*Note),
		tags:  make(map[string][]string),
	}

	result := &CreateNotebookResponse{Name: body.GetName()}
	response, err := json.Marshal(result)
	if err != nil {
		log.Fatalf("Unable to marshal response : %v", err)
	}
	w.Write(response)

}

// CreateNote takes a response body attempting to deserialise it to a
// CreateNotebookRequest object
func (n *NotebookRepo) CreateNote(w http.ResponseWriter, r *http.Request) {
	body := &CreateNoteRequest{}
	if err := jsonpb.Unmarshal(r.Body, body); err != nil {
		log.Fatalf("Unable to unmarshal message from request: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	notebook, ok := n.notebooks[body.GetNotebookName()]
	if !ok {
		errMsg := fmt.Sprintf("Notebook with name '%s' does not exist", body.GetNotebookName())
		log.Fatalf(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	if err := assertRequiredProperty(w, body.Title, "title"); err != nil {
		log.Fatalf(err.Error())
		return
	}
	if err := assertRequiredProperty(w, body.Body, "body"); err != nil {
		log.Fatalf(err.Error())
		return
	}

	id := uuid.New().String()
	createdPb := ptypes.TimestampNow()
	note := Note{
		Id:      id,
		Title:   body.Title,
		Body:    body.Body,
		Tags:    body.Tags,
		Created: createdPb,
	}

	notebook.notes[id] = &note
	for _, tag := range note.Tags {
		if !notebook.tags.tagHoldsNoteId(tag, id) {
			notebook.tags[tag] = append(notebook.tags[tag], id)
		}
	}

	result := &CreateNoteResponse{Id: id, Created: createdPb}
	response, err := json.Marshal(result)
	if err != nil {
		log.Fatalf("Unable to marshal response : %v", err)
	}
	w.Write(response)

}

// GetNote takes a response body attempting to deserialise it to a
// GetNotebookRequest object
func (n *NotebookRepo) GetNote(w http.ResponseWriter, r *http.Request) {
	body := &GetNoteRequest{}
	if err := jsonpb.Unmarshal(r.Body, body); err != nil {
		log.Fatalf("Unable to unmarshal message from request: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	notebook, ok := n.notebooks[body.GetNotebookName()]
	if !ok {
		errMsg := fmt.Sprintf("Notebook with name '%s' does not exist", body.GetNotebookName())
		log.Fatalf(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	note, ok := notebook.notes[body.GetId()]
	if !ok {
		errMsg := fmt.Sprintf("Note with id '%s' does not exist", body.GetId())
		log.Fatalf(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	result := &GetNoteResponse{Note: note}
	response, err := json.Marshal(result)
	if err != nil {
		log.Fatalf("Unable to marshal response : %v", err)
	}
	w.Write(response)

}

// UpdateNote takes a response body attempting to deserialise it to a
// UpdateNotebookRequest object
func (n *NotebookRepo) UpdateNote(w http.ResponseWriter, r *http.Request) {
	body := &UpdateNoteRequest{}
	if err := jsonpb.Unmarshal(r.Body, body); err != nil {
		log.Fatalf("Unable to unmarshal message from request: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	notebook, ok := n.notebooks[body.GetNotebookName()]
	if !ok {
		errMsg := fmt.Sprintf("Notebook with name '%s' does not exist", body.UpdateNotebookName())
		log.Fatalf(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	note, ok := notebook.notes[body.GetId()]
	if !ok {
		errMsg := fmt.Sprintf("Note with id '%s' does not exist", body.UpdateId())
		log.Fatalf(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}
	note.Tags

	result := &UpdateNoteResponse{Note: note}
	response, err := json.Marshal(result)
	if err != nil {
		log.Fatalf("Unable to marshal response : %v", err)
	}
	w.Write(response)

}

func (t tags) tagHoldsNoteId(tagName, noteID string) bool {
	noteIDs, ok := t[tagName]
	if !ok {
		return false
	}

	for _, id := range noteIDs {
		if id == noteID {
			return true
		}
	}
	return false
}

func tagsToAddAndRemove(oldTags, newTags []string) (add []string, remove []string) {
	for _, tag := range oldTags {
		for _, newTag := range newTags {
		}

	}
}

func assertRequiredProperty(w http.ResponseWriter, property interface{}, propertyName string) error {
	// For now only handle string casts
	if strProperty, ok := property.(string); ok {
		if strProperty == "" {
			errMsg := fmt.Sprintf(requiredProperty, propertyName)
			http.Error(w, errMsg, http.StatusBadRequest)
			return fmt.Errorf(errMsg)
		}
	}
	return nil
}

// r.HandleFunc("/notebook", CreateNotebook).Methods("POST")
// r.HandleFunc("/note", CreateNote).Methods("POST")
// r.HandleFunc("/notebook", GetNotebook).Methods("GET")
// r.HandleFunc("/note", GetNote).Methods("GET")
// r.HandleFunc("/note", UpdateNote).Methods("UPDATE")
// r.HandleFunc("/note", DeleteNote).Methods("DELETE")
