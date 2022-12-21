package main

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

type Entry struct {
	Text string
	Time time.Time
}

func (e Entry) Title() string       { return e.Time.String() }
func (e Entry) Description() string { return e.Text }
func (e Entry) FilterValue() string { return e.Text }

type Journal struct {
	Entries []Entry
}

func Save(e Entry) error {
	f, err := os.OpenFile("gournal.json", os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer f.Close()

	content, err := os.ReadFile("gournal.json")
	if err != nil {
		log.Fatal(err)
		return err
	}
	log.Printf("content: %s\n", content)
	var jsonContent []byte
	var journal Journal
	json.Unmarshal(content, &journal)
	log.Println(journal)
	if journal.Entries != nil {
		log.Println("break1")
		updatedJournal := Journal{append(journal.Entries, e)}
		jsonContent, err = json.Marshal(updatedJournal)
	} else {
		log.Println("break2")
		newJournal := Journal{[]Entry{e}}
		jsonContent, err = json.Marshal(newJournal)
	}
	if err != nil {
		log.Fatal(err)
		return err
	}

	if _, err = f.WriteAt(jsonContent, 0); err != nil {
		log.Fatal(err)
		return err
	}
	return err
}

func InitJournal() error {
	f, err := os.OpenFile("gournal.json", os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
		return err
	}

	defer f.Close()
	return nil
}

func GetEntries() (*Journal, error) {
	f, err := os.OpenFile("gournal.json", os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer f.Close()

	content, err := os.ReadFile("gournal.json")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	log.Printf("content: %s\n", content)
	var journal Journal
	err = json.Unmarshal(content, &journal)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &journal, err
}