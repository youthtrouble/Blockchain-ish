package main

import (
	"Blockchain_ish"
	"bufio"
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

var (
	bcServer chan []Blockchain_ish.Block
	Blockchainish []Blockchain_ish.Block
)



func replaceChain(newBlocks []Blockchain_ish.Block) {

	if len(newBlocks) > len(Blockchainish) {
		Blockchainish = newBlocks
	}

}

func handleotherConn(conn net.Conn) {
	defer conn.Close()

	io.WriteString(conn, "Enter the Temperature info [e.g: 37, Lagos]:")

	scanner := bufio.NewScanner(conn)

	// take in temperature info from stdin and add it to blockchain after conducting necessary validation
	go func() {
		for scanner.Scan() {
			tempInfo := scanner.Text()

			s := strings.Split(tempInfo, ",")
			Temperature := s[1]
			Location:= s[0]

			newBlock, err := Blockchain_ish.GenerateBlock(Blockchainish[len(Blockchainish)-1], Temperature, Location)
			if err != nil {
				log.Println(err)
				continue
			}
			if Blockchain_ish.IsBlockValid(newBlock, Blockchainish[len(Blockchainish)-1]) {
				newBlockchain := append(Blockchainish, newBlock)
				replaceChain(newBlockchain)
			}

			bcServer <- Blockchainish
			io.WriteString(conn, "\nEnter the Temperature info [e.g: 37, Lagos]:")
		}
	}()

	go func() {
		for {
			time.Sleep(30 * time.Second)
			newout, err := json.MarshalIndent(Blockchainish, "", "  ")
			if err != nil {
				log.Fatal(err)
			}
			io.WriteString(conn, string(newout))
			io.WriteString(conn, "\nEnter the Temperature info [e.g: 37, Lagos]:")
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

	bcServer = make(chan []Blockchain_ish.Block)

	server, err := net.Listen("tcp", ":"+os.Getenv("ADDRESS"))
	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()

	t := time.Now()
	genesisBlock := Blockchain_ish.Block{0, t.String(), "", "", "", ""}
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


