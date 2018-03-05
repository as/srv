package fs

import (
	"os"
	"path/filepath"
	"testing"
)

type testfile struct {
	name string
	data string
	rm   bool
	dir  bool
}

func (t testfile) Size() int64 {
	return int64(len(t.data))
}
func (t testfile) Bytes() []byte {
	return []byte(t.data)
}

var commondirs = []testfile{
	{".", "", false, true},
	{"..", "", false, true},
	{"../fs", "", false, true},
}
var createfiles = []testfile{
	{"fs.test.createfiles0", "hello world", true, false},
	{"fs.test.createfiles1", "h", true, false},
	{"fs.test.createfiles2", "hell545555555o world", true, false},
	{"fs.test.createfiles3", "helxxxxxxxxxxxxxlo world", true, false},
	{"fs.test.createfiles4", "", true, false},
}
var zerolenfile = testfile{"fs.test.zerolen", "", true, false}

func TestSizeOfZeroLengthFile(t *testing.T) {
	l := &Local{}
	testput(t, l, zerolenfile)
	defer clean(t, l, zerolenfile)
	teststat(t, l, zerolenfile)
}

func TestStatLocal(t *testing.T) {
	l := &Local{}
	testput(t, l, createfiles...)
	defer clean(t, l, createfiles...)
	teststat(t, l, createfiles...)
}
func TestStatRemote(t *testing.T) {
	srv := testServer(t, "tcp", "localhost:0")
	defer srv.Close()
	client := testClient(t, srv)

	teststat(t, client, commondirs...)
}

func TestStatRemoteDir(t *testing.T) {
	srv := testServer(t, "tcp", "localhost:0")
	client := testClient(t, srv)
	teststat(t, client, commondirs...)
}

func teststat(t *testing.T, l Fs, f ...testfile) {
	t.Helper()
	for i, f := range f {
		fi := teststat0(t, l, f)
		if fi == nil {
			return
		}
		if fi.Name() != filepath.Base(f.name) {
			t.Fatalf("pass %d: have name %q, want %q\n", i, fi.Name(), f.name)
		}
		if fi.Size() != f.Size() && !f.dir {
			t.Fatalf("pass %d: have size %d, want %d\n", i, fi.Size(), f.Size())
		}
		t.Log(fi)
	}
}

func teststat0(t *testing.T, fs Fs, f testfile) os.FileInfo {
	t.Helper()

	fi, err := fs.Stat(f.name)
	if err != nil {
		t.Fatalf("want no error, have %s", err)
		return nil
	}

	return fi
}
