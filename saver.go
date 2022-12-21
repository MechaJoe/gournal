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
	f1, err := os.OpenFile("gournal.json", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
		return err
	}

	defer f1.Close()
	return nil
}
