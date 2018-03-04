package fs

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"io/ioutil"
	"log"
	"net"
)

type Server struct {
	Local
	fd             net.Listener
	donec, donesrv chan bool
}

type client struct {
	conn net.Conn
	rx   chan []byte
	tx   chan []byte
}

func Serve(netw, addr string) (*Server, error) {
	fd, err := net.Listen(netw, addr)
	if err != nil {
		return nil, err
	}
	s := &Server{
		Local:   Local{},
		fd:      fd,
		donec:   make(chan bool, 1),
		donesrv: make(chan bool),
	}
	s.donec <- true

	//	go s.run()
	go func() {
		for {
			select {
			case <-s.donesrv:
				return
			default:
				conn, err := fd.Accept()
				if err != nil {
					log.Printf("accept: %s\n", err)
					continue
				}
				go s.handle(&client{conn, make(chan []byte), make(chan []byte)})
			}
		}
	}()
	return s, nil
}

func (s *Server) Close() error {
	select {
	case ok := <-s.donec:
		if ok {
			close(s.donec)
			close(s.donesrv)
		}
	default:
	}
	return nil
}

func (s *Server) handle(c *client) {
	bio := bufio.NewReader(c.conn)
	defer c.conn.Close()
	for {
		hdr := make([]byte, 3)
		_, err := io.ReadAtLeast(bio, hdr, len(hdr))
		if err != nil {
			log.Printf("invalid header: %s", err)
		}
		select {
		case <-s.donesrv:
			return
		default:
		}
		switch string(hdr) {
		case "Get", "Put", "Cmd":
			ln, err := bio.ReadSlice('\n')
			if err != nil {
				log.Printf("readslice: %s\n", err)
				break
			}
			ln = bytes.TrimSpace(ln)
			switch string(hdr) {
			case "Get":
				data, err := s.Local.Get(string(ln))
				if err != nil {
					log.Printf("get: %s\n", err)
				}

				err = binary.Write(c.conn, binary.BigEndian, int64(len(data)))
				if err != nil {
					log.Printf("get: write len: %s\n", err)
				}

				_, err = c.conn.Write(data)
				if err != nil {
					log.Printf("get: write: %s\n", err)
				}

			case "Put":
				n := int64(0)
				err = binary.Read(bio, binary.BigEndian, &n)

				if err != nil {
					log.Printf("put: %s\n", err)
				}

				if n < 0 {
					log.Printf("put: len<0\n")
					return
				}

				data, err := ioutil.ReadAll(io.LimitReader(bio, n))
				if err != nil {
					log.Printf("put: data read err: %s\n", err)
				}

				err = s.Local.Put(string(ln), data)
				if err != nil {
					log.Printf("put: local: %s", err)
					return
				}

			}
		default:
			log.Printf("bad cmd: %s", hdr)
		}
	}
}
