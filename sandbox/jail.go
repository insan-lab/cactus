// Copyright 2014 The Cactus Authors. All rights reserved.

package sandbox

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/hjr265/jail.go/jail"
)

type jailCell struct {
	cell *jail.Cell
}

// NewJailCell creates a Cell backed by the ptrace-based jail.
func NewJailCell() (Cell, error) {
	dir, err := ioutil.TempDir("", "cactus-jail-")
	if err != nil {
		return nil, err
	}
	return &jailCell{cell: &jail.Cell{Dir: dir}}, nil
}

func (j *jailCell) Command(name string, args ...string) *Cmd {
	jcmd := j.cell.Command(name, args...)
	return NewCmd(&jailRunner{cmd: jcmd})
}

func (j *jailCell) Create(name string) (*os.File, error) {
	return j.cell.Create(name)
}

func (j *jailCell) Dispose() error {
	return os.RemoveAll(j.cell.Dir)
}

type jailRunner struct {
	cmd *jail.Cmd
}

func (r *jailRunner) Run() error {
	return r.cmd.Run()
}

func (r *jailRunner) StdinPipe() (io.WriteCloser, error) {
	return r.cmd.StdinPipe()
}

func (r *jailRunner) StdoutPipe() (io.ReadCloser, error) {
	return r.cmd.StdoutPipe()
}

func (r *jailRunner) StderrPipe() (io.ReadCloser, error) {
	return r.cmd.StderrPipe()
}

func (r *jailRunner) Usages() Resources {
	return Resources{
		Cpu:    r.cmd.Usages.Cpu,
		Memory: r.cmd.Usages.Memory,
	}
}

// SetLimits is called before Run to propagate limits to the jail command.
func (r *jailRunner) SetLimits(res Resources) {
	r.cmd.Limits.Cpu = res.Cpu
	r.cmd.Limits.Memory = res.Memory
}
