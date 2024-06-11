package main

import (
	"context"
	"flag"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/jasonmf/cmdutils/fatalif"
	"github.com/pkg/browser"
	"golang.org/x/sync/errgroup"
)

var (
	testArtists = []string{
		"Foo Flippers",
		"Bap Company",
		"Corndogs of War",
		"Aggresor Dunx",
	}
	testTitles = []string{
		"Turn It Up",
		"One",
		"Wave of Nausea",
		"Fest of Vengeance",
	}

	Overlay []byte
)

var (
	fDemo = flag.Bool("demo", false, "demo mode")
)

func main() {
	flag.Parse()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	eg, ctx := errgroup.WithContext(ctx)

	d := &Decks{}

	if *fDemo {
		eg.Go(func() error {
			log.Print("starting demo mode")
			defer log.Print("exiting demo mode")
			timer := time.NewTicker(time.Second * 5)
			defer timer.Stop()
			deckID := 0
			for {
				select {
				case <-timer.C:
				case <-ctx.Done():
					return nil
				}
				artist := testArtists[rand.Intn(len(testArtists))]
				title := testTitles[rand.Intn(len(testTitles))]
				d.setDeck(deckID+1, Track{
					Artist: artist,
					Title:  title,
				})
				deckID = (deckID + 1) % 4
			}
		})
	} else {
		eg.Go(func() error {
			defer log.Print("exiting stagelinq")
			for {
				if err := ctx.Err(); err != nil {
					return err
				}
				if err := maim(ctx, d.Handle); err != nil {
					log.Print("error in stagelinq: ", err)
				}
				time.Sleep(time.Second)
			}
		})
	}

	homeDir, err := os.UserHomeDir()
	fatalif.Error(err, "determining home directory")
	overlayPath := filepath.Join(homeDir, "akslobs.html")

	eg.Go(func() error {
		defer log.Print("exiting webserver")
		mux := http.NewServeMux()

		mux.HandleFunc("/obs", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, overlayPath)
		})
		mux.HandleFunc("/decks", d.getDecks)
		NewUI(d, cancel, mux)
		listen := "localhost:8011"
		log.Print("starting webserver on ", listen)
		server := &http.Server{
			Addr:    listen,
			Handler: mux,
		}
		go func() {
			<-ctx.Done()
			shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()
			server.Shutdown(shutdownCtx)
		}()
		go func() {
			uiURL := "http://" + listen + "/ui"
			for {
				time.Sleep(time.Second / 4)
				log.Print("testing for server ready on ", uiURL)
				resp, err := http.Get(uiURL)
				if err != nil {
					continue
				}
				defer resp.Body.Close()
				io.Copy(io.Discard, resp.Body)
				break
			}
			log.Print("opening browser")
			if err := browser.OpenURL(uiURL); err != nil {
				log.Print("error opening browser: ", err)
				cancel()
			}
		}()
		return server.ListenAndServe()
	})
	eg.Wait()
}
