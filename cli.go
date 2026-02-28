package main

import (
	"flag"
	"fmt"
	"os"
)

type CLI struct {
	bc *Blockchain
}

func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println(`  addblock   -data "BLOCK_DATA"     Add a block to the blockchain`)
	fmt.Println("  printchain                     Print all the blocks of the blockchain")
	fmt.Println("  validate                       Validate PoW + links")
	fmt.Println("  height                         Print number of blocks (including genesis)")
	fmt.Println("  reset                          Delete DB and create a fresh chain")
	fmt.Println(`  getblock  -hash "BLOCK_HASH"    Print a single block by its hash (hex)`)
}

func (cli *CLI) Run() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}

	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	addBlockData := addBlockCmd.String("data", "", "Block data")

	getBlockCmd := flag.NewFlagSet("getblock", flag.ExitOnError)
	getBlockHash := getBlockCmd.String("hash", "", "Block hash in hex")

	switch os.Args[1] {
	case "addblock":
		_ = addBlockCmd.Parse(os.Args[2:])
		if *addBlockData == "" {
			cli.printUsage()
			os.Exit(1)
		}
		cli.bc.AddBlock(*addBlockData)
		fmt.Println("✅ Block added!")

	case "printchain":
		cli.printChain()

	case "validate":
		fmt.Printf("Chain valid: %v\n", cli.bc.IsValid())

	case "height":
		fmt.Printf("Height: %d\n", cli.bc.Height())

	case "reset":
		cli.bc.Close()
		ResetBlockchain()
		cli.bc = NewBlockchain()
		fmt.Println("✅ Reset complete (new genesis created).")

	case "getblock":
		_ = getBlockCmd.Parse(os.Args[2:])
		if *getBlockHash == "" {
			cli.printUsage()
			os.Exit(1)
		}
		cli.getBlock(*getBlockHash)

	default:
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) printChain() {
	it := cli.bc.Iterator()

	for {
		block := it.Next()
		printBlock(block)

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}

func (cli *CLI) getBlock(hashHex string) {
	block, raw, err := cli.bc.GetBlockHex(hashHex)
	if err != nil {
		fmt.Printf("❌ %v\n", err)
		fmt.Printf("hash: %x\n", raw)
		os.Exit(1)
	}

	printBlock(block)
}

func printBlock(block *Block) {
	fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
	fmt.Printf("Data: %s\n", block.Data)
	fmt.Printf("Hash: %x\n", block.Hash)
	fmt.Printf("Nonce: %d\n", block.Nonce)

	pow := NewProofOfWork(block)
	fmt.Printf("PoW valid: %v\n", pow.Validate())
	fmt.Println()
}
