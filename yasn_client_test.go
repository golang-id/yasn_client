package yasn_client

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

var (
	server *httptest.Server
	mux    *http.ServeMux
	client *Client
)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	client = NewClient(nil)
	client.BaseURL, _ = url.Parse(server.URL)
}

func tearDown() {
	server.Close()
}

func TestNewClient(t *testing.T) {
	testClient := NewClient(nil)
	if testClient.BaseURL.String() != DefaultBaseURL {
		t.Errorf("Client.BaseURL = '%s', wants '%s'", testClient.BaseURL, DefaultBaseURL)
	}
}

func TestGetNote(t *testing.T) {
	setup()
	defer tearDown()

	mux.HandleFunc("/notes/1", func(w http.ResponseWriter, r *http.Request) {
		expectMethod("GET", r, t)
		fmt.Fprint(w, `{"id": 1}`)
	})

	note, err := client.GetNote(1)
	if err != nil {
		t.Errorf("GetNote(1) returns error: %s", err)
	}

	expect := &Note{Id: 1}
	if !reflect.DeepEqual(note, expect) {
		t.Errorf("GetNote(1) returns %+v, wants %+v", note, expect)
	}
}

func TestAddNote(t *testing.T) {
	setup()
	defer tearDown()

	mux.HandleFunc("/notes", func(w http.ResponseWriter, r *http.Request) {
		expectMethod("POST", r, t)
		fmt.Fprint(w, `{"id": 1, "title": "test add note"}`)
	})

	newNote := &Note{
		Title: "test add note",
	}
	note, err := client.AddNote(newNote)
	if err != nil {
		t.Errorf("AddNote(%+v) returns error: %s", newNote, err)
	}

	expect := &Note{
		Id:    1,
		Title: "test add note",
	}
	if !reflect.DeepEqual(note, expect) {
		t.Errorf("AddNote(%+v) returns %+v, wants %+v", newNote, note, expect)
	}
}

func expectMethod(method string, r *http.Request, t *testing.T) {
	if method != r.Method {
		t.Errorf("request method = '%s', wants '%s'", r.Method, method)
	}
}
