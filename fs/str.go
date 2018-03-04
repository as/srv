package fs

import (
	"encoding/binary"
	"errors"
	"io"
	"io/ioutil"
)

var (
	ErrStrlen = errors.New("string too long")
)

// readstring reads a 64bit length-prefixed string
// in the form: [8]n [n]string
func readString(r io.Reader, max int64) (string, error) {
	n := int64(0)
	err := binary.Read(r, binary.BigEndian, &n)
	if err != nil {
		return "", err
	}

	if n > max {
		return "", ErrStrlen
	}
	data, err := ioutil.ReadAll(io.LimitReader(r, n))
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// writeString writes a 64bit length-prefixed string
// in the form: [8]n [n]string
func writeString(w io.Writer, s string) (err error) {
	if err = binary.Write(w, binary.BigEndian, int64(len(s))); err != nil {
		return err
	}
	n, err := w.Write([]byte(s))
	if err != nil {
		return err
	}
	if n != len(s) {
		panic("writeString: err != nil && short write")
	}
	return nil
}
