// Copyright 2014 The Cactus Authors. All rights reserved.

package sandbox

import (
	"io"
	"os"
	"time"
)

// Resources represents CPU and memory resource limits or usage.
type Resources struct {
	Cpu    time.Duration
	Memory uint64
}

// Cell represents an isolated environment for compiling and running code.
type Cell interface {
	Command(name string, args ...string) *Cmd
	Create(name string) (*os.File, error)
	Dispose() error
}

// Runner is the internal interface that sandbox implementations provide.
type Runner interface {
	Run() error
	StdinPipe() (io.WriteCloser, error)
	StdoutPipe() (io.ReadCloser, error)
	StderrPipe() (io.ReadCloser, error)
	SetLimits(Resources)
	Usages() Resources
}

// Cmd wraps a Runner with public Limits and Usages fields so that
// belt code can access them directly.
type Cmd struct {
	Limits Resources
	Usages Resources
	runner Runner
}

func NewCmd(runner Runner) *Cmd {
	return &Cmd{runner: runner}
}

func (c *Cmd) Run() error {
	c.runner.SetLimits(c.Limits)
	err := c.runner.Run()
	c.Usages = c.runner.Usages()
	return err
}

func (c *Cmd) StdinPipe() (io.WriteCloser, error) {
	return c.runner.StdinPipe()
}

func (c *Cmd) StdoutPipe() (io.ReadCloser, error) {
	return c.runner.StdoutPipe()
}

func (c *Cmd) StderrPipe() (io.ReadCloser, error) {
	return c.runner.StderrPipe()
}
