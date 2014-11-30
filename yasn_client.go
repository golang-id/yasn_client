package yasn_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	DefaultBaseURL = "http://localhost:8000/api/v1/"
)

type Note struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	ContentHTML string `json:"content_html"`
	Tags        []*Tag `json:"tags"`
}

type Tag struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type NotesAPIClient interface {
	GetNote(id int) (*Note, error)
	AddNote(note *Note) (*Note, error)

	// @todo GetAllNotes, EditNote, DeleteNote
}

type Client struct {
	// HTTP client used to communicate YASN API
	client *http.Client

	// Base URL for API request
	BaseURL *url.URL
}

func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	baseURL, _ := url.Parse(DefaultBaseURL)

	return &Client{
		client:  http.DefaultClient,
		BaseURL: baseURL,
	}
}

func (c *Client) GetNote(id int) (*Note, error) {
	ep := fmt.Sprintf("notes/%d", id)
	resp, err := c.request("GET", ep, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	note := new(Note)
	err = json.NewDecoder(resp.Body).Decode(note)
	if err != nil {
		return nil, err
	}
	return note, nil
}

func (c *Client) AddNote(note *Note) (*Note, error) {
	resp, err := c.request("POST", "notes", note)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	newNote := new(Note)
	err = json.NewDecoder(resp.Body).Decode(newNote)
	if err != nil {
		return nil, err
	}
	return newNote, nil
}

func (c *Client) request(method string, endpoint string, body interface{}) (*http.Response, error) {
	ep, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(ep)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if s := resp.StatusCode; s < 200 || s > 299 {
		// @todo returns error message if there's any
		return nil, fmt.Errorf("Unexpected status code %d", s)
	}

	return resp, nil
}
