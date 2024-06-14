package main

import (
	_ "embed"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

const (
	htmlName = "akslobs.html"
)

var (
	//go:embed akslobs.html
	overlayHTML []byte
)

func NewOverlayHandler(args []string) http.Handler {
	if htmlPath := fileFromArgs(args); htmlPath != "" {
		return serveFile(htmlPath)
	}
	if htmlPath := fileFromWD(); htmlPath != "" {
		return serveFile(htmlPath)
	}
	if htmlPath := fileFromHomeDir(); htmlPath != "" {
		return serveFile(htmlPath)
	}
	// give up, return the embedded file
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Content-Length", strconv.Itoa(len(overlayHTML)))
		w.Write(overlayHTML)
	})
}

func serveFile(path string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path)
	})
}

func fileViable(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func fileFromHomeDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	htmlPath := filepath.Join(homeDir, htmlName)
	if fileViable(htmlPath) {
		return htmlPath
	}
	return ""
}

func fileFromWD() string {
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}
	htmlPath := filepath.Join(wd, htmlName)
	if fileViable(htmlPath) {
		return htmlPath
	}
	return ""
}

func fileFromArgs(args []string) string {
	if len(args) < 1 {
		return ""
	}
	stat, err := os.Stat(args[0])
	if err != nil {
		return ""
	}
	if stat.IsDir() {
		// it's a directory. Maybe akslobs.html is there
		htmlPath := filepath.Join(args[0], htmlName)
		stat, err := os.Stat(htmlPath)
		if err != nil {
			return ""
		}
		if stat.IsDir() {
			return ""
		}
		return htmlPath
	}
	return args[0]
}
