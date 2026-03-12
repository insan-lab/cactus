// Copyright 2014 The Cactus Authors. All rights reserved.

package ui

import "embed"

//go:embed index.min.html
var indexHTML embed.FS

//go:embed assets
var assetsFS embed.FS
