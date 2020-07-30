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

// assertRequiredProperty is a convenience function to check if a request body
// has a non zero length required string field
func assertRequiredProperty(w http.ResponseWriter, property interface{}, propertyName string) error {
	// NOTE For now only handle string casts
	if strProperty, ok := property.(string); ok {
		if strProperty == "" {
			errMsg := fmt.Sprintf(requiredProperty, propertyName)
			http.Error(w, errMsg, http.StatusBadRequest)
			return fmt.Errorf(errMsg)
		}
	}
	return nil
}

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

// Notebook represents the notebook object holding a collection of Notes
// for the in memory NotebookRepo
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

// notesMeta is meant to be returned in a GetNotebookResponse.Notes slice
type notesMeta []*Note

// hasID checks if a noteID string is present in a notesMeta object
// as a Note object
func (m notesMeta) hasID(noteID string) bool {
	for _, note := range m {
		if note.Id == noteID {
			return true
		}
	}
	return false
}

// GetNotebook takes a response body attempting to deserialise it to a
func (n *NotebookRepo) GetNotebook(w http.ResponseWriter, r *http.Request) {
	body := &GetNotebookRequest{}
	if err := jsonpb.Unmarshal(r.Body, body); err != nil {
		log.Fatalf("Unable to unmarshal message from request: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	notebook, ok := n.notebooks[body.GetName()]
	if !ok {
		errMsg := fmt.Sprintf("Notebook with name '%s' does not exist", body.GetName())
		log.Fatalf(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	var notes notesMeta

	for _, note := range notebook.notes {
		if len(body.Tags) != 0 {
			for _, tag := range body.Tags {
				if notebook.tags.tagHoldsNoteID(tag, note.GetId()) && !notes.hasID(note.Id) {
					notes = append(notes, &Note{
						Id:      note.Id,
						Title:   note.Title,
						Body:    "",
						Tags:    note.Tags,
						Created: note.Created,
					})
				}
			}
		} else {
			notes = append(notes, &Note{
				Id:      note.Id,
				Title:   note.Title,
				Body:    "",
				Tags:    note.Tags,
				Created: note.Created,
			})
		}
	}

	result := &GetNotebookResponse{
		Name:  body.Name,
		Notes: notes,
	}

	response, err := json.Marshal(result)
	if err != nil {
		log.Fatalf("Unable to marshal response : %v", err)
	}
	w.Write(response)

}

// CreateNote takes a request body attempting to deserialise it to a
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
		if !notebook.tags.tagHoldsNoteID(tag, id) {
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
// GetNoteRequest object
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
	log.Println("UpdateNote: ", body)

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

	// compiles a list of tag names to determine which
	// ones should be added and removed by returning
	// two distinct string slices for addion and removal
	add, remove := tagsToAddAndRemove(note.Tags, body.Tags)

	// removal first to marginally reduce O(n) complexity
	// of addition
	for _, tagName := range remove {
		notebook.tags[tagName] = removeNoteID(notebook.tags[tagName], body.Id)
		// if no noteIDs are left in the tag slice
		if len(notebook.tags[tagName]) == 0 {
			delete(notebook.tags, tagName)
		}
	}

	for _, tagName := range add {
		notebook.tags[tagName] = addNoteID(notebook.tags[tagName], body.Id)
	}

	// update with everyting from new note expect timestamps
	note = &Note{
		Title:        body.Title,
		Body:         body.Body,
		Tags:         body.Tags,
		Created:      note.Created,
		LastModified: ptypes.TimestampNow(),
	}

	result := &UpdateNoteResponse{Note: note}
	response, err := json.Marshal(result)
	if err != nil {
		log.Fatalf("Unable to marshal response : %v", err)
	}
	w.Write(response)

}

// DeleteNote takes a response body attempting to deserialise it to a
// DeleteNotebookRequest object
func (n *NotebookRepo) DeleteNote(w http.ResponseWriter, r *http.Request) {
	body := &DeleteNoteRequest{}
	if err := jsonpb.Unmarshal(r.Body, body); err != nil {
		log.Fatalf("Unable to unmarshal message from request: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println("DeleteNote: ", body)

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

	// compiles a list of tag names to determine which
	// ones should be added and removed by returning
	// pass an empty slice for new tags
	_, remove := tagsToAddAndRemove(note.Tags, []string{""})

	for _, tagName := range remove {
		notebook.tags[tagName] = removeNoteID(notebook.tags[tagName], body.Id)
	}

	// once tags have been scrubbed, delete the key proper
	delete(notebook.notes, body.GetId())

	// update with everyting from new note expect timestamps
	note = &Note{
		Title:        note.Title,
		Body:         note.Body,
		Tags:         note.Tags,
		Created:      note.Created,
		LastModified: ptypes.TimestampNow(),
	}

	result := &DeleteNoteResponse{Note: note}
	response, err := json.Marshal(result)
	if err != nil {
		log.Fatalf("Unable to marshal response : %v", err)
	}
	w.Write(response)

}

// tagHoldsNoteId is used to determine whether a Notebook.tags entry
// holds a noteID reference
func (t tags) tagHoldsNoteID(tagName, noteID string) bool {
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

// tag operation is used by tagsToAddAndRemove to determine what tag
// references to insert or pop
type tagOperation struct {
	hasOld bool
	hasNew bool
}

// tagsToAddAndRemove is used when updating a note to determine
// what noteIDs to add and remove from Notebook.tags
func tagsToAddAndRemove(oldTags, newTags []string) (add []string, remove []string) {
	tagMap := make(map[string]*tagOperation)

	for _, tag := range oldTags {
		tagMap[tag] = &tagOperation{}
		tagMap[tag].hasOld = true
	}

	for _, tag := range newTags {
		if tagMap[tag] == nil {
			tagMap[tag] = &tagOperation{}
		}
		tagMap[tag].hasNew = true
	}

	for tagName, presence := range tagMap {
		// if found in new tags and missing from old *Note object
		if presence.hasNew && !presence.hasOld {
			add = append(add, tagName)
		}

		// if missing from new tagging and present in old *Note object
		if !presence.hasNew && presence.hasOld {
			remove = append(remove, tagName)
		}
	}
	return add, remove
}

// removeNoteID takes a slice of noteIDs and
// removes any reference to it in a returned slice
func removeNoteID(noteSlice []string, noteID string) (newSlice []string) {

	for _, id := range noteSlice {
		if id != noteID {
			newSlice = append(newSlice, id)
		}
	}
	return newSlice
}

// addNoteID takes a slice of noteIDs and adds noteID
// to it if it is not already present, returning said slice
func addNoteID(noteSlice []string, noteID string) []string {
	containsNoteID := false
	for _, id := range noteSlice {
		if id == noteID {
			containsNoteID = true
		}
	}
	if !containsNoteID {
		noteSlice = append(noteSlice, noteID)
	}
	return noteSlice
}
