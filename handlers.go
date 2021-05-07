package main

import (
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	"io"
	"net/http"
)

type Param struct {
	BPM int
}

func handleGetBlockchainish(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.MarshalIndent(Blockchainish, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

func handlePostBlockchainish(w http.ResponseWriter, r *http.Request) {
	var m Param

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	newBlock, err := generateBlock(Blockchainish[len(Blockchainish)-1], m.BPM)
	if err != nil {
		respondWithJSON(w, r, http.StatusInternalServerError, m)
		return
	}

	if isBlockValid(newBlock, Blockchainish[len(Blockchainish)-1]) {
		newBlockchain := append(Blockchainish, newBlock)
		replaceChain(newBlockchain)
		spew.Dump(Blockchainish)
	}

	respondWithJSON(w, r, http.StatusCreated, newBlock)
}

func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	w.WriteHeader(code)
	w.Write(response)
}
