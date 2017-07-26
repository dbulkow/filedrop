package main

import (
	"testing"
	"time"
)

func TestMkhash(t *testing.T) {
	hash := "0f3ea03220e341b00fd818f60ce91013fdba2bcea861ddb7d6d41a73a2cb4087"

	md := &MetaData{
		Files:   []File{{Name: "testfile"}},
		Created: time.Date(2001, 7, 8, 1, 15, 0, 0, time.UTC),
		Expire:  time.Date(2001, 7, 10, 1, 15, 0, 0, time.UTC),
	}

	md.MkHash()

	if md.Hash != hash {
		t.Fatalf("Expected hash %s got %s\n", hash, md.Hash)
	}
}
