package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
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


func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		t := time.Now()
		genesisBlock := Block{0, t.String(), 0, "", ""}
		spew.Dump(genesisBlock)
		Blockchainish = append(Blockchainish, genesisBlock)
	}()
	log.Fatal(run())

}


