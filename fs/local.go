package fs

import (
	"context"
	"io/ioutil"
	"os/exec"
)

type Local struct {
}

func (f *Local) Get(name string) (data []byte, err error) {
	return ioutil.ReadFile(name)
}
func (f *Local) Put(name string, data []byte) (err error) {
	return ioutil.WriteFile(name, data, 0666)
}
func (f *Local) Cmd(ctx context.Context, name string, arg ...string) (*exec.Cmd, error) {
	return exec.CommandContext(ctx, name, arg...), nil
}
