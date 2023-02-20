package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Block struct {
	Data         map[string]interface{}
	Hash         string
	PreviousHash string
	Timestamp    time.Time
	PoW          int
}

type Blockchain struct {
	GenesisBlock Block
	Chain        []Block
	Difficulty   int
}

func CreateBlockchain(difficulty int) Blockchain {
	genesisBlock := Block{
		Hash:      "0",
		Timestamp: time.Now(),
	}

	return Blockchain{
		GenesisBlock: genesisBlock,
		Chain:        []Block{genesisBlock},
		Difficulty:   difficulty,
	}
}

func (b Block) CalculateHash() string {
	data, _ := json.Marshal(b.Data)
	blockData := b.PreviousHash + strconv.Itoa(b.PoW) + string(data) + b.Timestamp.String()
	hash := sha256.Sum256([]byte(blockData))

	fmt.Println(fmt.Sprintf("Hash calculado: %x", hash))

	return fmt.Sprintf("%x", hash)
}

func (b *Block) Mine(difficulty int) {
	for !strings.HasPrefix(b.Hash, strings.Repeat("0", difficulty)) {
		b.PoW++
		b.Hash = b.CalculateHash()
	}
}

func (b *Blockchain) AddBlock(to, from string, amount float64) {
	data := map[string]interface{}{
		"to":     to,
		"from":   from,
		"amount": amount,
	}

	lastBlock := b.Chain[len(b.Chain)-1]

	newBlock := Block{
		Data:         data,
		PreviousHash: lastBlock.Hash,
		Timestamp:    time.Now(),
	}

	newBlock.Mine(b.Difficulty)

	fmt.Println("Novo bloco:")
	fmt.Println(newBlock)

	b.Chain = append(b.Chain, newBlock)
}

func (b Blockchain) Validate() bool {
	for i := range b.Chain[1:] {
		if b.Chain[i].Hash != "0" {
			previousBlock := b.Chain[i-1]
			currentBlock := b.Chain[i]

			if currentBlock.Hash != currentBlock.CalculateHash() {
				return false
			}

			if previousBlock.Hash != currentBlock.PreviousHash {
				return false
			}
		}
	}

	return true
}

func main() {
	// Create a new Blockchain
	blockchain := CreateBlockchain(3)

	// Create new transactions from Alice to Bob and John to Bob
	blockchain.AddBlock("Alice", "Bob", 5)
	blockchain.AddBlock("John", "Bob", 2)

	// Validate Blockchain, expect True
	fmt.Println(fmt.Sprintf("Is valid blockchain: %t", blockchain.Validate()))
}
