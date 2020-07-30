package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var newNotebook = &NotebookRepo{
	notebooks: map[string]Notebook{
		"new_notebook": {
			notes: map[string]*Note{"note_1": {
				Id:    "id_1",
				Title: "title_1",
				Body:  "body_1",
				Tags:  []string{"tag_1"},
			}},
			tags: tags{
				"tag_1": []string{"note_1"},
			},
		},
	},
}

// meant to be used as a GetNotebookResponse
var newNotebookNoBody = &NotebookRepo{
	notebooks: map[string]Notebook{
		"new_notebook": {
			notes: map[string]*Note{"note_1": {
				Id:    "id_1",
				Title: "title_1",
				Tags:  []string{"tag_1"},
			}},
		},
	},
}

func TestCreateNotebook(t *testing.T) {
	repo := NewNotebookRepo()
	t.Run("Success", func(t *testing.T) {
		req := httptest.NewRequest(
			"POST",
			"/notebook",
			bytes.NewBufferString(`{"name": "new_notebook"}`),
		)
		w := httptest.NewRecorder()
		repo.CreateNotebook(w, req)

		resp := w.Result()

		assert.Equal(t, 200, resp.StatusCode)
	})
	t.Run("NoteBookExists/Error", func(t *testing.T) {
		repo = newNotebook
		req := httptest.NewRequest(
			"POST",
			"/notebook",
			bytes.NewBufferString(`{"name": "new_notebook"}`),
		)
		w := httptest.NewRecorder()
		repo.CreateNotebook(w, req)

		resp := w.Result()

		assert.Equal(t, 409, resp.StatusCode)
	})
}

func TestGetNotebook(t *testing.T) {
	repo := newNotebook
	t.Run("NoFilter/Success", func(t *testing.T) {
		req := httptest.NewRequest(
			"GET",
			"/notebook",
			bytes.NewBufferString(`{"name": "new_notebook"}`),
		)
		w := httptest.NewRecorder()
		repo.GetNotebook(w, req)

		resp := w.Result()

		body, _ := ioutil.ReadAll(resp.Body)

		expected := &GetNotebookResponse{
			Name: "new_notebook",
			Notes: []*Note{
				{
					Id:    "id_1",
					Title: "title_1",
					Tags:  []string{"tag_1"},
				},
			},
		}
		expectedJson, _ := json.Marshal(expected)
		assert.Equal(t, string(expectedJson), string(body))
		assert.Equal(t, 200, resp.StatusCode)
	})
	t.Run("UnfulfilledFilter", func(t *testing.T) {
		req := httptest.NewRequest(
			"GET",
			"/notebook",
			bytes.NewBufferString(`{"name": "new_notebook", "tags": ["missing_filter"]}`),
		)
		w := httptest.NewRecorder()
		repo.GetNotebook(w, req)

		resp := w.Result()

		body, _ := ioutil.ReadAll(resp.Body)

		expected := &GetNotebookResponse{
			Name: "new_notebook",
		}
		expectedJson, _ := json.Marshal(expected)
		assert.Equal(t, string(expectedJson), string(body))
		assert.Equal(t, 200, resp.StatusCode)
	})
}
func TestCreateNote(t *testing.T) {
	repo := newNotebook
	req := httptest.NewRequest(
		"POST",
		"/note",
		bytes.NewBufferString(`
{
  "notebook_name": "new_notebook",
  "title": "title_2",
  "body": "body_2",
  "tags": [
    "tag_2"
  ]
}
	`),
	)
	w := httptest.NewRecorder()
	repo.CreateNote(w, req)

	resp := w.Result()

	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, 200, resp.StatusCode)
	assert.NotNil(t, body)
}

// Ran out of time for test coverage  ¯\_(ツ)_/¯
