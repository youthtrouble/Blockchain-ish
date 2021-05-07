package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"time"
)

type Block struct {
	Index		int
	Timestamp	string
	BPM 		int
	Hash		string
	PrevHash	string
}

var Blockchainish []Block

func run() error {
	mux := spunNewRouter()
	httpAddr := os.Getenv("ADDRESS")
	log.Println("Listening on ", os.Getenv("ADDRESS"))
	s := &http.Server{
		Addr:           ":" + httpAddr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}


func spunNewRouter() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/", handleGetBlockchainish).Methods("GET")
	router.HandleFunc("/", handlePostBlockchainish).Methods("POST")
	return router
}
