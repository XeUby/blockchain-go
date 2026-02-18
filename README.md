# Minimal Blockchain in Go

[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![bbolt](https://img.shields.io/badge/storage-bbolt-blue)](https://github.com/etcd-io/bbolt)

A clean, single-node **blockchain implementation** written in pure Go — designed as an educational and portfolio project.

Demonstrates core blockchain concepts:

- Cryptographic hashing & immutability  
- Proof-of-Work (PoW) mining  
- Persistent storage with BoltDB (bbolt)  
- Chain validation & integrity checks  
- Simple but powerful CLI interface  

Perfect for learning blockchain fundamentals, practicing systems programming in Go, and showcasing backend engineering skills during interviews.

## ✨ Features

- SHA-256 Proof-of-Work with adjustable difficulty  
- Persistent storage using BoltDB (embedded key-value store)  
- Full chain integrity validation (hashes + PoW + links)  
- CLI commands: add block, print chain, validate, height, reset  
- Genesis block creation on first run or reset  
- Minimal dependencies (only `bbolt` + standard library)  
- Clean, layered & well-commented code  

## 🚀 Quick Start

# Clone and run
git clone https://github.com/XeUby/blockchain-go.git
cd blockchain-go

# Reset database and create genesis block
go run . reset

# Add your first block
go run . addblock -data "Send 10 BTC to Walter"

# See the whole chain
go run . printchain

# Validate integrity
go run . validate

# Show current height
go run . height
## 📋 CLI Commands

| Command                              | Description                                      |
|--------------------------------------|--------------------------------------------------|
| `go run . reset`                     | Delete DB and create fresh genesis block         |
| `go run . addblock -data "..."`      | Mine and add a new block with data               |
| `go run . printchain`                | Print all blocks (from genesis to tip)           |
| `go run . validate`                  | Check full chain integrity                       |
| `go run . height`                    | Show current blockchain height                   |

## 🔗 Block Structure

Each block contains the following fields:

```go
type Block struct {
    Timestamp     int64
    Data          []byte
    PrevBlockHash []byte
    Hash          []byte        // SHA-256 hash of the block header
    Nonce         int64
}
```
- Timestamp — creation time 
- Data — arbitrary payload (e.g. "Send 10 BTC to Walter")
- PrevBlockHash — hash of the previous block (links blocks into a chain)
- Hash — current block's hash (computed during mining)
- Nonce — value found during Proof-of-Work

Blocks are cryptographically linked: changing any data in a past block invalidates all subsequent blocks.

## ⛏️ Proof-of-Work
The mining process finds a nonce such that:
textSHA256(PrevBlockHash + Data + Timestamp + Nonce) < Target

Target is a very small number derived from the targetBits constant.
Lower targetBits → smaller target → exponentially harder mining.
Mining difficulty can be adjusted by changing targetBits in pow.go.

This simulates Bitcoin-style Proof-of-Work and demonstrates computational work required for immutability.
## 💾 Persistence Layer
Blocks are stored using bbolt (modern, reliable fork of BoltDB):

Bucket: "blocks"
Key = block hash → Value = serialized block (Gob-encoded)
Special key: "lh" (last hash) → hash of the current chain tip

This design ensures:

Fast lookup by hash
Automatic persistence across restarts
Atomic updates when adding new blocks

## 📂 Project Structure

```go
textblockchain-go/
├── main.go          # Entry point & CLI wiring
├── cli.go           # Command-line parser and handlers
├── blockchain.go    # Core blockchain logic + DB operations
├── block.go         # Block struct + serialization / deserialization
├── pow.go           # Proof-of-Work mining & validation logic
├── go.mod           # Module definition
└── README.md
```
All files are kept minimal and focused — easy to read and understand in one sitting.
## 🎯 Design Goals

- Clarity first — code is simple, well-commented, and easy to follow
- Correctness — full chain validation (hashes + PoW + links)
- Minimalism — only essential dependencies (bbolt + stdlib)
- Educational value — clearly shows how immutability and consensus basics work
- Interview/portfolio friendly — compact project that demonstrates systems thinking

## ⚠️ Limitations
This is an educational implementation, not production software:

- Single-node only (no peer-to-peer networking)
- No real transactions, wallets, signatures or UTXO model
- No Merkle trees or advanced features
- No difficulty retargeting algorithm
- No REST API (only CLI for simplicity)

## 🛠 Tech Stack

- Go 1.26.0
- bbolt — embedded, crash-safe key-value store
- crypto/sha256 — standard library hashing

## 🚀 Possible Improvements 

- Implement basic transactions and simple wallet (ECDSA keys)
- Add Merkle Tree for transaction integrity
- Introduce P2P networking (e.g. via libp2p or custom TCP)
- Dynamic difficulty adjustment every N blocks
- REST/JSON API for easier integration
- Unit & integration tests (especially chain validation)
- Docker container + docker-compose
- Visualization of chain / mining process

Built with ❤️ by Boris Chugin
