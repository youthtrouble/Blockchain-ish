package Blockchain_ish

//This file primarily holds the Block logic and instance
import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

type Block struct {
	Index		int
	Timestamp	string
	TEMPERATURE string
	Location	string  //considered using time.Location but that only provides the timezone so better not to :)
	Hash		string
	PrevHash	string
}

func calculateHash(block Block) string {

	record := string(block.Index) + block.Timestamp + string(block.TEMPERATURE) + block.Location + block.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)

}

func GenerateBlock(oldBlock Block, Temperature string, Location string) (Block, error) {

	var newBlock Block

	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.TEMPERATURE = Temperature
	newBlock.Location = Location
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = calculateHash(newBlock)

	return newBlock, nil
}

func IsBlockValid(newBlock, oldBlock Block) bool{
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
