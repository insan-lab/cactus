// Copyright 2014 The Cactus Authors. All rights reserved.

package belt

import (
	"io"
	"time"

	"github.com/FurqanSoftware/cactus/sandbox"
)

type StackC struct{}

func (s *StackC) Build(cell sandbox.Cell, source io.Reader) (*sandbox.Cmd, error) {
	f, err := cell.Create("source.c")
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(f, source)
	if err != nil {
		return nil, err
	}
	err = f.Close()
	if err != nil {
		return nil, err
	}

	cmd := cell.Command("gcc", "source.c")
	cmd.Limits.Cpu = 16 * time.Second
	cmd.Limits.Memory = 1 << 30

	return cmd, nil
}

func (s *StackC) Run(cell sandbox.Cell) *sandbox.Cmd {
	return cell.Command("./a.out")
}
