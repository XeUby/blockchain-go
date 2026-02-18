# Blockchain in Go

A minimal single-node blockchain implementation written in Go.

This project demonstrates the core mechanics behind blockchain systems, including block structure, cryptographic hashing, Proof-of-Work mining, persistent storage, and validation.  
It is designed as an educational and portfolio project to showcase backend engineering, data integrity principles, and systems-level reasoning in Go.

---

## Overview

The blockchain consists of sequentially linked blocks.  
Each block:

- Stores arbitrary data
- References the previous block via its hash
- Is mined using a Proof-of-Work algorithm
- Is persisted locally using BoltDB

Tampering with any historical block invalidates the chain due to hash dependencies and PoW verification.

---

## Features

- SHA-256 based Proof-of-Work consensus mechanism
- Adjustable difficulty via `targetBits`
- Persistent storage using BoltDB (bbolt)
- CLI interface for interaction
- Full chain validation (hash + PoW + link integrity)
- Blockchain height inspection
- Database reset functionality

---

# CLI Usage

### Reset the blockchain (create a fresh genesis block)
go run . reset

Add a new block
go run . addblock -data "Send 10 BTC to Walter"

Print the entire blockchain
go run . printchain

Validate the blockchain
go run . validate

Show blockchain height
go run . height

Block Structure

---

## Block Structure. Each block contains:

Timestamp
Data
PrevBlockHash
Hash
Nonce

Blocks are cryptographically linked through PrevBlockHash, forming an immutable chain.

---

## Proof-of-Work

Mining is performed by iterating over nonce values until the following condition is satisfied:
SHA256(block_data) < target

The target is derived from a difficulty parameter (targetBits).
Increasing difficulty raises the expected mining time.

Persistence Layer

Blocks are stored in BoltDB as:
hash → serialized block

The latest block hash (the chain tip) is stored under the key:

lh → last hash

This guarantees persistence across application restarts.

---

## Project Structure
block.go        → Block definition and serialization
pow.go          → Proof-of-Work implementation
blockchain.go   → Core blockchain logic and persistence
cli.go          → Command-line interface
main.go         → Application entry point

---

## Design Goals

Demonstrate blockchain immutability
Emphasize correctness and validation
Keep the architecture minimal and readable
Avoid unnecessary external dependencies
Focus on core mechanics rather than networking
Limitations
Single-node implementation (no distributed consensus)
No transaction model (data stored as raw payload)
No peer-to-peer networking
Educational project, not production-ready
---
## Tech Stack

Go
BoltDB (bbolt)
SHA-256 (crypto/sha256)

---

## Author
Boris Chugin
GitHub: https://github.com/XeUby
