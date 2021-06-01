package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
)


type Block struct {
	Index		int
	Timestamp	string
	BPM 		int
	Hash		string
	PrevHash	string
}

type Param struct {
	BPM int
}
var bcServer chan []Block


var Blockchainish []Block

func calculateHash(block Block) string {

	record := string(block.Index) + block.Timestamp + string(block.BPM) + block.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)

}

func generateBlock(oldBlock Block, BPM int) (Block, error) {

	var newBlock Block

	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.BPM = BPM
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = calculateHash(newBlock)

	return newBlock, nil
}

func isBlockValid(newBlock, oldBlock Block) bool{
	if oldBlock.Index + 1 != newBlock.Index {
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}

	if calculateHash(newBlock) != newBlock.Hash {
		return false
	}

	return true
}

func replaceChain(newBlocks []Block) {

	if len(newBlocks) > len(Blockchainish) {
		Blockchainish = newBlocks
	}

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

func handleotherConn(conn net.Conn) {
	defer conn.Close()

	io.WriteString(conn, "Enter a new BPM:")

	scanner := bufio.NewScanner(conn)

	// take in BPM from stdin and add it to blockchain after conducting necessary validation
	go func() {
		for scanner.Scan() {
			bpm, err := strconv.Atoi(scanner.Text())
			if err != nil {
				log.Printf("%v not a number: %v", scanner.Text(), err)
				continue
			}
			newBlock, err := generateBlock(Blockchainish[len(Blockchainish)-1], bpm)
			if err != nil {
				log.Println(err)
				continue
			}
			if isBlockValid(newBlock, Blockchainish[len(Blockchainish)-1]) {
				newBlockchain := append(Blockchainish, newBlock)
				replaceChain(newBlockchain)
			}

			bcServer <- Blockchainish
			io.WriteString(conn, "\nEnter a new BPM:")
		}
	}()
}


func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	bcServer = make(chan []Block)

	server, err := net.Listen("tcp", ":"+os.Getenv("ADDRESS"))
	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()

	t := time.Now()
	genesisBlock := Block{0, t.String(), 0, "", ""}
	spew.Dump(genesisBlock)
	Blockchainish = append(Blockchainish, genesisBlock)

	log.Fatal(run())


	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleotherConn(conn)
	}

}


