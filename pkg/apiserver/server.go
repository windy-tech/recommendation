package apiserver

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"

	language "cloud.google.com/go/language/apiv1"
	"github.com/gorilla/mux"
	"github.com/windy-tech/recommendation/pkg/apiserver/handler"
)

type APIServer struct {
	Router     *mux.Router
	DB         *sql.DB
	LangClient *language.Client
}

func (a *APIServer) setLangClient() {
	log.Println("Preparing ML server")
	ctx := context.Background()

	// Creates a client.
	client, err := language.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	a.LangClient = client
}

func (a *APIServer) setRouters() {
	a.Get("/books", a.handleRequest(handler.GetAllBooks))
	a.Post("/books", a.handleRequest(handler.CreateBook))
}

// Get wraps the router for GET method
func (a *APIServer) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("GET")
}

// Post wraps the router for POST method
func (a *APIServer) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("POST")
}

// Put wraps the router for PUT method
func (a *APIServer) Put(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("PUT")
}

// Delete wraps the router for DELETE method
func (a *APIServer) Delete(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("DELETE")
}

// Run the aPIAPIServer on it's router
func (a *APIServer) Run() {
	port := "8080"
	if p := os.Getenv("APP_PORT"); p != "" {
		port = p
	}
	a.Router = mux.NewRouter()
	a.setRouters()
	a.setLangClient()
	log.Fatal(http.ListenAndServe("localhost:"+port, a.Router))
}

type RequestHandlerFunction func(db *sql.DB, lang *language.Client, w http.ResponseWriter, r *http.Request)

func (a *APIServer) handleRequest(handler RequestHandlerFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(a.DB, a.LangClient, w, r)
	}
}
