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

## CLI Usage

Reset the blockchain (creates fresh genesis block):

go run . reset

Add a new block:

go run . addblock -data "Send 10 BTC to Alice"

Print the blockchain:

go run . printchain

Validate the blockchain:

go run . validate

Show blockchain height:

go run . height
How It Works
Block Structure

Each block contains:

Timestamp

Data

PrevBlockHash

Hash

Nonce

Blocks are cryptographically linked through PrevBlockHash.

Proof-of-Work

Mining is performed by iterating over nonce values until:

SHA256(block_data) < target

The target is derived from a difficulty parameter (targetBits).
Higher difficulty increases expected mining time.

Persistence Layer

Blocks are stored in BoltDB as:

hash → serialized block

The latest block hash is stored under the key:

lh → last hash (tip)

This ensures persistence across application restarts.

Project Structure
block.go        → Block definition and serialization
pow.go          → Proof-of-Work implementation
blockchain.go   → Core blockchain logic and persistence
cli.go          → Command-line interface
main.go         → Application entry point
Design Goals

Demonstrate blockchain immutability

Emphasize correctness and validation

Keep architecture minimal and readable

Avoid external dependencies beyond storage

Focus on core mechanics rather than networking

Limitations

Single-node (no networking or distributed consensus)

No transaction model (data is stored as raw payload)

No peer-to-peer communication

Educational implementation, not production-ready

Tech Stack

Go

BoltDB (bbolt)

SHA-256 (crypto/sha256)

Author

Boris Chugin
GitHub: https://github.com/XeUby
