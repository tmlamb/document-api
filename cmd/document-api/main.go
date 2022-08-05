package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"

	"github.com/tmlamb/document-api/cmd/document-api/handler"
)

func router() http.Handler {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.Use(middleware)
	myRouter.PathPrefix("/").Methods("OPTIONS")
	myRouter.HandleFunc("/accounts", handler.ReturnAllAccounts).Methods("GET")
	myRouter.HandleFunc("/accounts", handler.CreateNewAccount).Methods("POST")
	myRouter.HandleFunc("/accounts/{accountId}", handler.ReturnAccount).Methods("GET")
	myRouter.HandleFunc("/accounts/{accountId}", handler.UpdateAccount).Methods("PUT")
	myRouter.HandleFunc("/accounts/{accountId}", handler.DeleteAccount).Methods("DELETE")
	myRouter.HandleFunc("/accounts/{accountId}/documents", handler.UploadDocument).Methods("POST")
	myRouter.HandleFunc("/accounts/{accountId}/documents/index", handler.ReturnDocumentsIndex).Methods("GET")
	myRouter.HandleFunc("/accounts/{accountId}/documents/{id}", handler.DownloadDocument).Methods("GET")
	myRouter.HandleFunc("/accounts/{accountId}/documents/{id}", handler.UpdateDocument).Methods("PUT")
	myRouter.HandleFunc("/accounts/{accountId}/documents/{id}", handler.DeleteDocument).Methods("DELETE")
	return myRouter
}

func middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("middleware", r.URL)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "*")

		if (*r).Method == "OPTIONS" {
			log.Println("OPTIONS preflight")
			return
		}

		h.ServeHTTP(w, r)
	})
}

func main() {
	conn, err := pgx.Connect(context.Background(),
		fmt.Sprintf("postgres://%s:%s",
			os.Getenv("PGHOST"),
			os.Getenv("PGPORT")))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close(context.Background())

	log.Fatal(http.ListenAndServe(":8080", router()))
}
