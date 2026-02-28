package main

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"

	bolt "go.etcd.io/bbolt"
)

const (
	dbFile       = "blockchain.db"
	blocksBucket = "blocks"
	lastHashKey  = "lh"
)

var ErrBlockNotFound = errors.New("block not found")

type Blockchain struct {
	tip []byte
	db  *bolt.DB
}

type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

func NewBlockchain() *Blockchain {
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	var tip []byte

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		if b == nil {
			// Create bucket + genesis
			b, err = tx.CreateBucket([]byte(blocksBucket))
			if err != nil {
				return err
			}

			genesis := NewGenesisBlock()
			if err := b.Put(genesis.Hash, genesis.Serialize()); err != nil {
				return err
			}
			if err := b.Put([]byte(lastHashKey), genesis.Hash); err != nil {
				return err
			}
			tip = genesis.Hash
		} else {
			tip = b.Get([]byte(lastHashKey))
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return &Blockchain{tip: tip, db: db}
}

func (bc *Blockchain) Close() {
	_ = bc.db.Close()
}

func (bc *Blockchain) AddBlock(data string) {
	var lastHash []byte

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte(lastHashKey))
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	newBlock := NewBlock(data, lastHash)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		if err := b.Put(newBlock.Hash, newBlock.Serialize()); err != nil {
			return err
		}
		if err := b.Put([]byte(lastHashKey), newBlock.Hash); err != nil {
			return err
		}

		bc.tip = newBlock.Hash
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{currentHash: bc.tip, db: bc.db}
}

func (it *BlockchainIterator) Next() *Block {
	var block *Block

	err := it.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encoded := b.Get(it.currentHash)
		if encoded == nil {
			return ErrBlockNotFound
		}
		block = DeserializeBlock(encoded)
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	it.currentHash = block.PrevBlockHash
	return block
}

func (bc *Blockchain) Height() int {
	it := bc.Iterator()
	height := 0

	for {
		block := it.Next()
		height++

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
	return height
}

// IsValid checks:
// 1) Each block's PoW is valid
// 2) Each block correctly references the previous block hash
func (bc *Blockchain) IsValid() bool {
	it := bc.Iterator()

	for {
		block := it.Next()

		// PoW must be valid for current block
		pow := NewProofOfWork(block)
		if !pow.Validate() {
			return false
		}

		// Genesis has no previous
		if len(block.PrevBlockHash) == 0 {
			break
		}

		// Previous block must exist and its hash must match PrevBlockHash
		prev, err := bc.GetBlock(block.PrevBlockHash)
		if err != nil {
			return false
		}
		if !bytes.Equal(prev.Hash, block.PrevBlockHash) {
			return false
		}
	}
	return true
}

// GetBlock fetches a block by raw hash bytes.
func (bc *Blockchain) GetBlock(hash []byte) (*Block, error) {
	var block *Block

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encoded := b.Get(hash)
		if encoded == nil {
			return ErrBlockNotFound
		}
		block = DeserializeBlock(encoded)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return block, nil
}

// GetBlockHex fetches a block by hex string hash (e.g. "0000abc...").
// Supports optional "0x" prefix.
func (bc *Blockchain) GetBlockHex(hashHex string) (*Block, []byte, error) {
	s := hashHex
	if len(s) >= 2 && (s[:2] == "0x" || s[:2] == "0X") {
		s = s[2:]
	}

	raw, err := hex.DecodeString(s)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid hash hex: %w", err)
	}

	block, err := bc.GetBlock(raw)
	if err != nil {
		return nil, raw, err
	}
	return block, raw, nil
}

// Reset deletes DB file and creates a fresh chain with new genesis.
func ResetBlockchain() {
	_ = os.Remove(dbFile)
	// Recreate by calling NewBlockchain once
	bc := NewBlockchain()
	bc.Close()
}
