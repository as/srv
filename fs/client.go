package fs

import (
	"context"
	"encoding/binary"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os/exec"
)

type Client struct {
	conn net.Conn
}

func Dial(netw, addr string) (*Client, error) {
	conn, err := net.Dial(netw, addr)
	if err != nil {
		log.Printf("dial: %s\n", err)
		return nil, err
	}
	return &Client{conn}, nil
}

func (f *Client) Get(name string) (data []byte, err error) {
	f.conn.Write([]byte("Get" + name + "\n"))
	n := int64(0)
	err = binary.Read(f.conn, binary.BigEndian, &n)
	if err != nil {
		log.Printf("get: write len: %s\n", err)
	}
	return ioutil.ReadAll(io.LimitReader(f.conn, n))
}
func (f *Client) Put(name string, data []byte) (err error) {
	f.conn.Write([]byte("Put" + name + "\n"))
	err = binary.Write(f.conn, binary.BigEndian, int64(len(data)))
	if err != nil {
		log.Printf("put: write len: %s\n", err)
	}

	_, err = f.conn.Write(data)
	if err != nil {
		log.Printf("get: write: %s\n", err)
	}
	return err
}
func (f *Client) Cmd(ctx context.Context, name string, arg ...string) (*exec.Cmd, error) {
	return nil, nil
}
func (f *Client) Close() error {
	return f.conn.Close()
}
