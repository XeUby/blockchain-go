package main

import (
	"bytes"
	"testing"
)

func TestSerializeDeserializeBlock(t *testing.T) {
	orig := NewBlock("hello", []byte("prevhash"))

	encoded := orig.Serialize()
	decoded := DeserializeBlock(encoded)

	if orig.Timestamp != decoded.Timestamp {
		t.Fatalf("timestamp mismatch: %d != %d", orig.Timestamp, decoded.Timestamp)
	}
	if !bytes.Equal(orig.Data, decoded.Data) {
		t.Fatalf("data mismatch: %q != %q", orig.Data, decoded.Data)
	}
	if !bytes.Equal(orig.PrevBlockHash, decoded.PrevBlockHash) {
		t.Fatalf("prev hash mismatch")
	}
	if !bytes.Equal(orig.Hash, decoded.Hash) {
		t.Fatalf("hash mismatch")
	}
	if orig.Nonce != decoded.Nonce {
		t.Fatalf("nonce mismatch: %d != %d", orig.Nonce, decoded.Nonce)
	}
}
