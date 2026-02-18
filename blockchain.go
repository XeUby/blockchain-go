package main

import (
	"bytes"
	"fmt"
	"log"

	bolt "go.etcd.io/bbolt"
)

const (
	dbFile       = "blockchain.db"
	blocksBucket = "blocks"
	lastHashKey  = "lh"
)

type Blockchain struct {
	tip []byte
	db  *bolt.DB
}

type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

func NewBlockchain() *Blockchain {
	db, err := bolt.Open(dbFile, 0o600, nil)
	if err != nil {
		log.Panic(err)
	}

	var tip []byte

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		if b == nil {
			genesis := NewGenesisBlock()

			nb, e := tx.CreateBucket([]byte(blocksBucket))
			if e != nil {
				return e
			}

			if e := nb.Put(genesis.Hash, genesis.Serialize()); e != nil {
				return e
			}
			if e := nb.Put([]byte(lastHashKey), genesis.Hash); e != nil {
				return e
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
		if b == nil {
			return fmt.Errorf("bucket %q not found", blocksBucket)
		}

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
		block = DeserializeBlock(encoded)
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	it.currentHash = block.PrevBlockHash
	return block
}

func (bc *Blockchain) IsValid() bool {
	it := bc.Iterator()

	for {
		block := it.Next()

		pow := NewProofOfWork(block)
		if !pow.Validate() {
			return false
		}

		if len(block.PrevBlockHash) == 0 {
			return true
		}

		var prevExists bool
		_ = bc.db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(blocksBucket))
			prevExists = b.Get(block.PrevBlockHash) != nil
			return nil
		})
		if !prevExists {
			return false
		}

		if bytes.Equal(block.Hash, block.PrevBlockHash) {
			return false
		}
	}
}

// Height returns nuber of blocks including genesis.
func (bc *Blockchain) Height() int {
	it := bc.Iterator()
	count := 0

	for {
		block := it.Next()
		count++
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
	return count
}
