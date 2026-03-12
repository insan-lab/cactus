// Copyright 2014 The Cactus Authors. All rights reserved.

package ui

import (
	"io/fs"
	"net/http"
)

func ServeAsset(w http.ResponseWriter, r *http.Request) {
	sub, err := fs.Sub(assetsFS, "assets")
	catch(err)
	http.FileServer(http.FS(sub)).ServeHTTP(w, r)
}
