package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net"
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
var (
	bcServer chan []Block
	Blockchainish []Block
)

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

	go func() {
		for {
			time.Sleep(30 * time.Second)
			newout, err := json.Marshal(Blockchainish)
			if err != nil {
				log.Fatal(err)
			}
			io.WriteString(conn, string(newout))
			io.WriteString(conn, "\nEnter a new BPM:")
		}
	}()

	for _ = range bcServer {
		spew.Dump(Blockchainish)
	}
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

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleotherConn(conn)
	}

}


