package main

import (
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	mtr_go "github.com/untoreh/mtr-go"
	"github.com/untoreh/mtr-go/tools"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Print("loading services...")
	m := mtr_go.New(map[string]interface{}{
		"services": []string{
			"bing",
			"google",
			"yandex",
			"frengly",
			"multillect",
			"promt",
			"systran",

			// "convey",
			// "sdl",
		},
	})
	tools.Cache.Save()

	// Routes consist of a path and a handler function.
	r := mux.NewRouter()
	r.Handle("/", handlers.CompressHandler(&mtr_go.MtrGet{m})).Methods("GET")
	r.Handle("/", handlers.CompressHandler(&mtr_go.MtrPost{m})).Methods("POST")
	r.Handle("/multi", handlers.CompressHandler(&mtr_go.MtrPostMulti{m})).Methods("POST")

	// Bind to a port and pass our router in
	log.Print("starting server...")
	server := &http.Server{
		Addr:    ":8001",
		Handler: r,
	}
	log.Print(server.ListenAndServe())
}
