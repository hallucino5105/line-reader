package main

import (
	"fmt"
	"testing"
)

var testfilePath = "./testdata/bufio.txt"

func createLineReader(t *testing.T) *LineReader {
	mr, err := NewLineReader(testfilePath)
	if err != nil {
		t.Fatal(err)
	}
	return mr
}

func TestLineReader1(t *testing.T) {
	var num int64 = 30
	var skip int64 = 1

	mr := createLineReader(t)
	defer mr.Close()

	buf, err := mr.ReadLineMultiN(num, skip)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(buf))
}

func TestLineReader2(t *testing.T) {
	mr := createLineReader(t)
	defer mr.Close()

	buf, err := mr.ReadAll()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(buf))
}
