package main

import "testing"

func TestPoWValidForNewBlock(t *testing.T) {
	b := NewBlock("test data", []byte("prevhash"))
	pow := NewProofOfWork(b)

	if !pow.Validate() {
		t.Fatalf("expected PoW to be valid for a newly mined block")
	}
}

func TestPoWInvalidAfterDataTamper(t *testing.T) {
	b := NewBlock("original", []byte("prevhash"))
	pow := NewProofOfWork(b)

	if !pow.Validate() {
		t.Fatalf("expected PoW valid before tampering")
	}

	// tamper with data
	b.Data = []byte("tampered")

	// PoW should become invalid because Validate recomputes hash over block fields
	if pow.Validate() {
		t.Fatalf("expected PoW to be invalid after data tampering")
	}
}
