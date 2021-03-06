package line_reader

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"
	"github.com/stoicperlman/fls"
)

type EOFFileError struct {
}

func (p EOFFileError) Error() string {
	return io.EOF.Error()
}

type LineReader struct {
	Filepath string
	LineNum  int64

	fd         *fls.File
	buf        []byte
	delimiters []byte
}

func NewLineReader(filepath string) (*LineReader, error) {
	p := &LineReader{
		Filepath:   filepath,
		buf:        make([]byte, 32*1024),
		delimiters: []byte{'\n'},
	}

	p.Open()

	return p, nil
}

func (p *LineReader) Open() error {
	fd, err := fls.OpenFile(p.Filepath, os.O_RDONLY, 0600)
	if err != nil {
		return errors.Wrap(err, "can not open file")
	}
	p.fd = fd

	num, err := p.lineCounter()
	if err != nil {
		return errors.Wrap(err, "failed line count")
	}
	p.LineNum = num

	return nil
}

func (p *LineReader) Close() {
	if p.fd != nil {
		p.fd.Close()
		p.fd = nil
	}
}

func (p *LineReader) ReadLineMulti(num int64) ([]byte, error) {
	return p.ReadLineMultiN(num, 0)
}

func (p *LineReader) ReadLineMultiN(num int64, skip int64) ([]byte, error) {
	if p.fd == nil {
		return nil, fmt.Errorf("file not open")
	}

	p.seek(0, io.SeekStart)

	if skip != 0 {
		_, err := p.seek(skip, io.SeekStart)
		if err != nil {
			if err == io.EOF {
				if p.LineNum <= num {
					return p.ReadAll()
				} else {
					return p.ReadLineMultiN(num, p.LineNum-num)
				}
			} else {
				return nil, errors.Wrap(err, "failed line seek")
			}
		}
	}

	return p.read(num)
}

func (p *LineReader) ReadAll() ([]byte, error) {
	if p.fd == nil {
		return nil, fmt.Errorf("file not open")
	}
	p.seek(0, io.SeekStart)
	return p.read(p.LineNum)
}

func (p *LineReader) seek(line int64, whence int) (int64, error) {
	if line == 0 && whence == io.SeekStart {
		return p.fd.Seek(0, io.SeekStart)
	}
	return p.fd.SeekLine(line, whence)
}

func (p *LineReader) read(num int64) ([]byte, error) {
	var count int64 = 0
	var buf = []byte{}

	scanner := bufio.NewScanner(p.fd)

	for scanner.Scan() {
		line := append(scanner.Bytes(), '\n')
		buf = append(buf, line...)

		count += 1
		if num <= count {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, errors.Wrap(err, "failed read line")
	}

	return buf, nil
}

func (p *LineReader) lineCounter() (int64, error) {
	p.seek(0, io.SeekStart)

	var count int64 = 0
	for {
		c, err := p.fd.Read(p.buf)
		count += int64(bytes.Count(p.buf[:c], p.delimiters))

		switch {
		case err == io.EOF:
			return count, nil
		case err != nil:
			return count, err
		}
	}
}
