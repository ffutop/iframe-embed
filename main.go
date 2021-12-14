package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
)

var iframeHTML string = "<!doctype html><html><frameset rows=\"100%%\"><frame src=\"%s\"><noframes><a href=\"%s\">Click here</a></noframes></frameset></html>\n"

type iframeHandler struct {
	embedHost string
}

var handler iframeHandler

func (h iframeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(fmt.Sprintf(iframeHTML, h.embedHost+path, h.embedHost+path)))
	w.WriteHeader(http.StatusOK)
}

func main() {
	flag.StringVar(&handler.embedHost, "embedHost", "https://workspace.easyv.cloud", "iframe embed host")

	router := mux.NewRouter()

	router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	router.PathPrefix("/").Handler(handler)

	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
