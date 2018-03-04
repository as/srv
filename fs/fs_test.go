package fs

import (
	"os"
	"testing"
)

func teststat0(t *testing.T, fs Fs, name string, errok bool) os.FileInfo {
	t.Helper()
	fi, err := fs.Stat(name)
	if err != nil {
		if errok {
			return nil
		}
		t.Logf("want no error, have %s", err)
		t.Fail()
		return nil
	}

	if errok {
		t.Logf("have no error, want error")
		t.Fail()
	}

	return fi
}

func TestStatLocal(t *testing.T) {
	teststat(t, &Local{})
}
func TestStatRemote(t *testing.T) {
	srv := testServer(t, "tcp", "localhost:0")
	client := testClient(t, srv)
	teststat(t, client)
}

func teststat(t *testing.T, l Fs) {
	t.Helper()
	name := "fs.test.stat"
	data := "stat test"
	size := int64(len(data))
	testput(t, l, name, []byte(data), false)

	fi := teststat0(t, l, name, false)
	if fi == nil {
		return
	}
	if fi.Name() != name {
		t.Logf("have name %q, want %q\n", fi.Name(), name)
		t.Fail()
	}
	if fi.Size() != size {
		t.Logf("have size %d, want %d\n", fi.Size(), size)
		t.Fail()
	}
	t.Log(fi)
}

func testput(t *testing.T, fs Fs, name string, data []byte, rm bool) {
	t.Helper()
	err := fs.Put(name, data)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	if !rm {
		return
	}
	err = os.Remove(name)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
}

func testget(t *testing.T, fs Fs, name string, rm bool) (data []byte) {
	t.Helper()

	have, err := fs.Get(name)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	if rm {
		if err = os.Remove(name); err != nil {
			t.Log(err)
			t.Fail()
		}
	}

	return have
}

func TestFSPut(t *testing.T) {
	l := &Local{}
	testput(t, l, "fs.test.write", []byte("hello world"), true)
}

func TestFSGet(t *testing.T) {
	l := &Local{}
	name := "fs.test.get"
	want := "take me to your leader"

	testput(t, l, name, []byte(want), false)
	have := testget(t, l, name, true)

	if string(have) != want {
		t.Logf("data mismatch: have %q, want %q\n", have, want)
		t.Fail()
	}
}

func testServer(t *testing.T, netw, addr string) *Server {
	srv, err := Serve("tcp", "localhost:0")
	if err != nil {
		t.Logf("serve: %s\n", err)
		t.Fail()
	}
	return srv
}

func testClient(t *testing.T, srv *Server) *Client {
	t.Helper()
	addr := srv.fd.Addr()

	client, err := Dial(addr.Network(), addr.String())
	if err != nil {
		t.Logf("dial: %s\n", err)
		t.FailNow()
		return nil
	}
	return client
}

func TestServerClient(t *testing.T) {
	srv := testServer(t, "tcp", "localhost:0")
	client := testClient(t, srv)
	defer srv.Close()

	name := "fs.net.test.get"
	want := "take me to your leader"

	testput(t, client, name, []byte(want), false)
	have := testget(t, client, name, true)

	if string(have) != want {
		t.Logf("data mismatch: have %q, want %q\n", have, want)
		t.Fail()
	}
}
