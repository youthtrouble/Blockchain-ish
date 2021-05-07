package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
	"log"
	"time"
)





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


