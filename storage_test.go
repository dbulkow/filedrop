package main

import (
	"testing"
	"time"
)

func TestMkhash(t *testing.T) {
	hash := "d14a547dd8bf5aea00f53c72818a1b722d00d5a60eb9e6ad03296536dec0fd16"

	md := &MetaData{
		Filename: "testfile",
		Created:  time.Date(2001, 7, 8, 1, 15, 0, 0, time.UTC),
		Expire:   time.Date(2001, 7, 10, 1, 15, 0, 0, time.UTC),
	}

	md.mkhash()

	if md.Hash != hash {
		t.Fatalf("Expected hash %s got %s\n", hash, md.Hash)
	}
}
