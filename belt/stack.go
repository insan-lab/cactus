// Copyright 2014 The Cactus Authors. All rights reserved.

package belt

import (
	"io"

	"github.com/FurqanSoftware/cactus/sandbox"
)

type Stack interface {
	Build(cell sandbox.Cell, source io.Reader) (*sandbox.Cmd, error)
	Run(cell sandbox.Cell) *sandbox.Cmd
}

var Stacks = map[string]Stack{
	"c":   &StackC{},
	"cpp": &StackCpp{},
}
