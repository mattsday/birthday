package main

import (
	"cloud.google.com/go/firestore"
	"context"
	firebase "firebase.google.com/go/v4"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var fs *firestore.Client

const defaultPort = 8080

type pageData struct {
	Title string
}

func main() {
	// Connect to Firestore
	ctx := context.Background()
	fb := &firebase.Config{}
	app, err := firebase.NewApp(ctx, fb)
	if err != nil {
		log.Fatalln(err)
	}
	fs, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if port == 0 || err != nil {
		log.Printf("Port not set or invalid - setting to %v", defaultPort)
		port = defaultPort
	}

	r := mux.NewRouter()
	s := http.StripPrefix("/html/static/", http.FileServer(http.Dir("./resources/html/")))
	r.Use(webHandler)
	hc := r.PathPrefix("/health").Subrouter()
	hc.HandleFunc("", health).Methods("GET")
	hc.HandleFunc("/liveness", health).Methods("GET")
	hc.HandleFunc("/readiness", health).Methods("GET")
	site := r.PathPrefix("/").Subrouter()
	site.PathPrefix("/html/static/").Handler(s)
	site.HandleFunc("/", home)
	//site.HandleFunc("/logged-out", loginPage).Methods(http.MethodGet)
	log.Println("Startup Complete, listening on port", port)
	log.Fatal(http.ListenAndServe(":"+fmt.Sprintf("%v", port), r))
}

func webHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.RequestURI, "/v1") || strings.HasPrefix(r.RequestURI, "/api") || strings.HasPrefix(r.RequestURI, "/health") || strings.HasPrefix(r.RequestURI, "/actuator") {
			w.Header().Add("Content-Type", "application/json")
		} else if strings.HasSuffix(r.RequestURI, ".css") {
			w.Header().Add("Content-Type", "text/css")
		} else if strings.HasSuffix(r.RequestURI, ".jpg") {
			w.Header().Add("Content-Type", "image/jpg")
		} else if strings.HasPrefix(r.RequestURI, "/static") {
			w.Header().Add("Content-Type", "text/plain")
		} else {
			w.Header().Add("Content-Type", "text/html; charset=utf8")
		}
		next.ServeHTTP(w, r)
	})
}

func health(w http.ResponseWriter, r *http.Request) {
	// A very simple health check.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(w, `{"alive": true}`)
}

func home(w http.ResponseWriter, r *http.Request) {
	pd := pageData{
		Title: "Hello World",
	}

	t := template.Must(template.New("index.tmpl").ParseFiles("resources/templates/index.tmpl", "resources/templates/header.tmpl", "resources/templates/nav.tmpl", "resources/templates/footer.tmpl"))
	//w.Write([]byte("Hello world"))
	err := t.Execute(w, pd)
	if err != nil {
		log.Printf("Error rendering page: %v\n", err)
		return
	}
}
