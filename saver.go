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
	User    string
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
	var jsonContent []byte
	var journal Journal
	json.Unmarshal(content, &journal)
	if journal.Entries != nil {
		updatedJournal := Journal{append(journal.Entries, e), journal.User}
		jsonContent, err = json.Marshal(updatedJournal)
	} else {
		newJournal := Journal{[]Entry{e}, journal.User}
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

func UpdateName(name string) error {
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
	var jsonContent []byte
	var journal Journal
	json.Unmarshal(content, &journal)
	updatedJournal := Journal{journal.Entries, name}
	jsonContent, err = json.Marshal(updatedJournal)
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

func GetJournal() (*Journal, error) {
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
	var journal Journal
	err = json.Unmarshal(content, &journal)
	if err != nil {
		return nil, err
	}
	return &journal, err
}
