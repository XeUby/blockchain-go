# Blockchain in Go

A minimal blockchain implementation written in Go.

This project demonstrates the core concepts behind a blockchain system, including block structure, hashing, Proof-of-Work, persistent storage, and a command-line interface. It is intended as an educational and portfolio project to showcase backend and systems programming skills in Go.

## Features

- Proof-of-Work consensus mechanism using SHA-256
- Persistent blockchain storage with BoltDB (bbolt)
- Command-line interface (CLI) for interacting with the blockchain
- Chain validation to verify block integrity and links

## Usage

Add a new block to the blockchain:
```bash
go run . addblock -data "Send 10 BTC to Alice"
Print all blocks in the blockchain:

go run . printchain
Validate the blockchain:

go run . validate
Project Structure
block.go — block definition and serialization logic

pow.go — Proof-of-Work implementation

blockchain.go — blockchain core logic and persistent storage

cli.go — command-line interface commands

main.go — application entry point

Tech Stack
Go

BoltDB (bbolt)

Author
Boris Chugin
GitHub: https://github.com/XeUby