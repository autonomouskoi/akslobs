package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/icedream/go-stagelinq"
)

type Track struct {
	Artist string `json:"artist"`
	Title  string `json:"title"`
}

type Decks struct {
	Deck1 Track `json:"deck1,omitempty"`
	Deck2 Track `json:"deck2,omitempty"`
	Deck3 Track `json:"deck3,omitempty"`
	Deck4 Track `json:"deck4,omitempty"`

	Hide [4]bool `json:"hide"`

	// in-progress
	deck1, deck2, deck3, deck4 Track
}

func (d *Decks) Handle(state *stagelinq.State) {
	log.Print("handling ", *state)
	// nothing / Engine / Deck / State
	nameFields := strings.Split(state.Name, "/")
	var currentDeck, inProgressDeck *Track
	log.Printf("2 %s, 3 %s", nameFields[2], nameFields[3])
	switch nameFields[2] {
	case "Deck1":
		currentDeck = &d.Deck1
		inProgressDeck = &d.deck1
	case "Deck2":
		currentDeck = &d.Deck2
		inProgressDeck = &d.deck2
	case "Deck3":
		currentDeck = &d.Deck3
		inProgressDeck = &d.deck3
	case "Deck4":
		currentDeck = &d.Deck4
		inProgressDeck = &d.deck4
	default:
		return
	}
	if nameFields[3] == "PlayState" {
		v, ok := state.Value["state"].(bool)
		if !ok || !v {
			return
		}
		*currentDeck = *inProgressDeck
		return
	}
	if len(nameFields) < 5 {
		return
	}
	if nameFields[3] != "Track" {
		return
	}
	switch nameFields[4] {
	case "ArtistName":
		v, ok := state.Value["string"].(string)
		if !ok {
			return
		}
		inProgressDeck.Artist = v
	case "SongName":
		v, ok := state.Value["string"].(string)
		if !ok {
			return
		}
		inProgressDeck.Title = v
	default:
		return
	}
}

func (d *Decks) getDecks(w http.ResponseWriter, _ *http.Request) {
	b, err := json.Marshal(d)
	if err != nil {
		http.Error(w, "marshalling deck data: "+err.Error(), http.StatusInternalServerError)
		log.Print("error marshalling deck data: ", err)
		return
	}
	w.Header().Set("Content-Type", "text/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(b)))
	if _, err := w.Write(b); err != nil {
		log.Print("error writing deck data: ", err)
	}
}

func (d *Decks) setDeck(id int, t Track) error {
	switch id {
	case 1:
		d.Deck1 = t
	case 2:
		d.Deck2 = t
	case 3:
		d.Deck3 = t
	case 4:
		d.Deck4 = t
	default:
		return errors.New("invalid deck ID")
	}
	return nil
}
