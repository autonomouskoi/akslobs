package main

import (
	_ "embed"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type UI struct {
	decks      *Decks
	keepalives *int
}

//go:embed ui.html
var uiHTML []byte

func NewUI(decks *Decks, cancel func(), mux *http.ServeMux) {
	ui := &UI{
		decks: decks,
	}
	mux.HandleFunc("/ui", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Content-Length", strconv.Itoa(len(uiHTML)))
		w.Write(uiHTML)
	})
	mux.HandleFunc("/ka", func(w http.ResponseWriter, r *http.Request) {
		zero := 0
		if ui.keepalives == nil {
			ui.keepalives = &zero
			return
		}
		*ui.keepalives = 0
	})
	mux.HandleFunc("/hide_deck", func(w http.ResponseWriter, r *http.Request) {
		deck := r.URL.Query().Get("deck")
		if len(deck) != 5 {
			return
		}
		if !strings.HasPrefix(deck, "deck") {
			return
		}
		id, err := strconv.Atoi(deck[4:])
		if err != nil {
			return
		}
		if id < 1 || id > 4 {
			return
		}
		decks.Hide[id-1] = r.URL.Query().Get("hide") == "true"
	})
	go func() {
		for {
			time.Sleep(time.Second * 1)
			if ui.keepalives == nil {
				continue
			}
			if *ui.keepalives > 5 {
				cancel()
				return
			}
			*ui.keepalives += 1
		}
	}()
}
