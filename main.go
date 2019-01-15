// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var (
	indexTmpl = template.Must(
		template.ParseFiles(filepath.Join("templates", "index.html")),
	)
)

func main() {
	http.HandleFunc("/", indexHandler)

	staticImage := http.StripPrefix("/image", http.FileServer(http.Dir("image")))
	staticCSS := http.StripPrefix("/css", http.FileServer(http.Dir("css")))
	staticJS := http.StripPrefix("/js", http.FileServer(http.Dir("js")))
	http.Handle("/image/", staticImage)
	http.Handle("/css/", staticCSS)
	http.Handle("/js/", staticJS)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

// indexHandler uses a template to create an index.html.
func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	type indexData struct {
		Logo            string
		Style           string
		JQuery          string
		BootstrapBundle string
		RequestTime     string
	}
	data := indexData{
		Logo:            "/image/logo.png",
		Style:           "/css/bootstrap.min.css",
		JQuery:          "/js/jquery.slim.min.js",
		BootstrapBundle: "/js/bootstrap.bundle.min.js",
		RequestTime: time.Now().Format(time.RFC822),
	}
	if err := indexTmpl.Execute(w, data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}