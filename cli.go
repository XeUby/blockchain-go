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
	fmt.Println("  addblock   -data \"BLOCK_DATA\"    Add a block to the blockchain")
	fmt.Println("  printchain                    Print all the blocks of the blockchain")
	fmt.Println("  validate                      Validate PoW + links")
	fmt.Println("  height                        Print number of blocks (including genesis)")
	fmt.Println("  reset                         Delete DB and create a fresh chain")
	fmt.Println()
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) addBlock(data string) {
	cli.bc.AddBlock(data)
	fmt.Println("✅ Block added!")
}

func (cli *CLI) printChain() {
	it := cli.bc.Iterator()

	for {
		block := it.Next()

		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Printf("Nonce: %d\n", block.Nonce)

		pow := NewProofOfWork(block)
		fmt.Printf("PoW valid: %v\n", pow.Validate())
		fmt.Println()

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}

func (cli *CLI) validateChain() {
	fmt.Println("Chain valid:", cli.bc.IsValid())
}

func (cli *CLI) height() {
	fmt.Println("Height:", cli.bc.Height())
}

func (cli *CLI) reset() {
	cli.bc.Close()
	_ = os.Remove(dbFile) // ignore if file doesn't exist

	cli.bc = NewBlockchain()
	fmt.Println("✅ Reset complete (new genesis created).")
}

func (cli *CLI) Run() {
	cli.validateArgs()

	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	validateCmd := flag.NewFlagSet("validate", flag.ExitOnError)
	heightCmd := flag.NewFlagSet("height", flag.ExitOnError)
	resetCmd := flag.NewFlagSet("reset", flag.ExitOnError)

	addBlockData := addBlockCmd.String("data", "", "Block data")

	switch os.Args[1] {
	case "addblock":
		_ = addBlockCmd.Parse(os.Args[2:])
	case "printchain":
		_ = printChainCmd.Parse(os.Args[2:])
	case "validate":
		_ = validateCmd.Parse(os.Args[2:])
	case "height":
		_ = heightCmd.Parse(os.Args[2:])
	case "reset":
		_ = resetCmd.Parse(os.Args[2:])
	default:
		cli.printUsage()
		os.Exit(1)
	}

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			fmt.Println("❌ -data is required")
			addBlockCmd.Usage()
			os.Exit(1)
		}
		cli.addBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

	if validateCmd.Parsed() {
		cli.validateChain()
	}

	if heightCmd.Parsed() {
		cli.height()
	}

	if resetCmd.Parsed() {
		cli.reset()
	}
}
