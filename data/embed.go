// Copyright 2014 The Cactus Authors. All rights reserved.

package data

import "embed"

//go:embed db-init.sql
var dbInitSQL embed.FS
