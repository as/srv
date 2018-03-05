package fs

import (
	"os"
	"sync"
	"testing"
)

func TestFSPut(t *testing.T) {
	l := &Local{}
	testput(t, l, createfiles...)
	clean(t, l, createfiles...)
}

func TestFSGet(t *testing.T) {
	l := &Local{}
	testput(t, l, createfiles...)
	defer clean(t, l, createfiles...)
	testget(t, l, createfiles...)
}

func TestServerClient(t *testing.T) {
	srv := testServer(t, "tcp", "localhost:0")
	defer srv.Close()

	var wg sync.WaitGroup
	for i := range createfiles {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			c := testClient(t, srv)
			defer c.Close()
			testput(t, c, createfiles[i])
			defer clean(t, c, createfiles[i])
			testget(t, c, createfiles[i])
			clean(t, c, createfiles[i])
		}()
	}
	wg.Wait()
}

func clean(t *testing.T, fs Fs, f ...testfile) {
	t.Helper()
	for _, f := range f {
		if !f.rm {
			t.Logf("clean: refusing rm: %s", f.name)
			continue
		}
		if err := os.Remove(f.name); err != nil {
			t.Log(err)
		}
	}
}

func testput(t *testing.T, fs Fs, f ...testfile) {
	t.Helper()
	for i, f := range f {
		if !f.rm {
			t.Fatalf("pass %d: put: %q: %s:", i, f.name, "cant write over a file with rm==false")
		}
		if err := fs.Put(f.name, f.Bytes()); err != nil {
			t.Fatalf("pass %d: put: %q: %s:", i, f.name, err)
		}
	}
}

func testget(t *testing.T, fs Fs, f ...testfile) {
	t.Helper()

	for i, f := range f {
		want := f.data
		have, err := fs.Get(f.name)
		if err != nil {
			t.Fatalf("pass %d: get %s: %s", i, f.name, err)
		}
		if have := string(have); have != want {
			t.Fatalf("get: contents differ: %s\n\thave: %q\n\twant%q\n", f.name, have, want)
		}
	}
	return
}

func testServer(t *testing.T, netw, addr string) *Server {
	srv, err := Serve("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("serve: %s\n", err)
	}
	return srv
}

func testClient(t *testing.T, srv *Server) *Client {
	t.Helper()
	addr := srv.fd.Addr()

	client, err := Dial(addr.Network(), addr.String())
	if err != nil {
		t.Fatalf("dial: %s\n", err)
		return nil
	}
	return client
}
