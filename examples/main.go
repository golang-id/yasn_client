package main

import (
	"fmt"
	"os"

	"github.com/golang-id/yasn_client"
)

func main() {
	client := yasn_client.NewClient(nil)

	// Get a note.
	note, err := client.GetNote(1)
	if err != nil {
		fmt.Printf("error: %+v\n", err)
		os.Exit(1)
	}
	printNote(note)

	// Add new note.
	newNote := &yasn_client.Note{
		Title:   "Test Add Note",
		Content: "Test Note",
	}
	retNote, err := client.AddNote(newNote)
	if err != nil {
		fmt.Printf("error: %+v\n", err)
		os.Exit(1)
	}
	printNote(retNote)
}

func printNote(n *yasn_client.Note) {
	fmt.Printf("Note.Id: %d\n", n.Id)
	fmt.Printf("Note.Title: %s\n", n.Title)
	fmt.Printf("Note.Content: %s\n", n.Content)
	fmt.Printf("Note.ContentHTML: %s\n", n.ContentHTML)
}
