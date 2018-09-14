package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
)

const difficulty = 1

type Block struct {
	Index      int    // Position of the data recorded in the blockchain
	Timestamp  string // The time the data is written
	Msg        string // Anything
	Hash       string // SHA256 identifier representing this data record
	PrevHash   string // SHA256 identifier of the previous record in the chain
	Nonce      string // Increments during mining process to facilitate hashing
	Difficulty int    // Constant defining the number of 0's in the leading hash
}

var Blockchain []Block

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		t := time.Now()
		genesisBlock := Block{0, t.String(), "GENESIS BLOCK", "", "", "", difficulty}
		// Prints structs to the console for useful debugging
		spew.Dump(genesisBlock)
		Blockchain = append(Blockchain, genesisBlock)
	}()
	log.Fatal(run())

}

func calculateHash(block Block) string {
	// Makes a hash out of the contents of the block
	record := string(block.Index) + block.Timestamp + string(block.Msg) + block.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func generateBlock(oldBlock Block, Msg string) (Block, error) {

	var newBlock Block

	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.Msg = Msg
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Difficulty = difficulty

	// Proof of work!
	for i := 0; ; i++ {
		hex := fmt.Sprintf("%x", i)
		newBlock.Nonce = hex
		if !isHashValid(calculateHash(newBlock), newBlock.Difficulty) {
			fmt.Println(calculateHash(newBlock), " do more work!")
			// Work
			time.Sleep(time.Second)
			continue
		} else {
			fmt.Println(calculateHash(newBlock), " mined!")
			newBlock.Hash = calculateHash(newBlock)
			break
		}

	}
	return newBlock, nil
}

func isHashValid(hash string, difficulty int) bool {
	prefix := strings.Repeat("0", difficulty)
	return strings.HasPrefix(hash, prefix)
}

// Makes sure there's no funny business
func isBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
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

// Take the longer chain if there's a dispute. It's more up to date.
func replaceChain(newBlocks []Block) {
	if len(newBlocks) > len(Blockchain) {
		Blockchain = newBlocks
	}
}
