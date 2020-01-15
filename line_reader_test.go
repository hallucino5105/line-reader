package line_reader

import (
	"fmt"
	"io"
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

func TestLineReader3(t *testing.T) {
	mr := createLineReader(t)
	defer mr.Close()

	buf1, err := mr.ReadLineMultiN(10, 5)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Print(string(buf1))
	fmt.Println("===")

	buf2, err := mr.ReadLineMultiN(10, 5)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Print(string(buf2))
}

func TestLineReader4(t *testing.T) {
	mr := createLineReader(t)
	defer mr.Close()

	mr.seek(0, io.SeekStart)
	buf, err := mr.read(50)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Print(string(buf))

	pos, err := mr.seek(0, io.SeekCurrent)
	fmt.Println(pos, err)
}
