// Copyright 2014 The Cactus Authors. All rights reserved.

package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"

	"github.com/hjr265/go-zrsc/zrsc"
	"github.com/pelletier/go-toml"

	"github.com/FurqanSoftware/cactus/belt"
)

var cfg *toml.Tree

func main() {
	runtime.GOMAXPROCS((runtime.NumCPU() + 1) / 2)

	_, err := os.Stat("config.toml")
	if os.IsNotExist(err) {
		f2, err := zrsc.Open("cmd/cactus/config-sample.toml")
		catch(err)

		f, err := os.Create("config.toml")
		_, err = io.Copy(f, f2)
		catch(err)

		err = f2.Close()
		catch(err)
		err = f.Close()
		catch(err)
	}

	cfg, err = toml.LoadFile("config.toml")
	catch(err)

	go func() {
		addr, ok := cfg.Get("core.addr").(string)
		if !ok {
			log.Fatal("Missing core.addr in config.toml")
		}

		log.Printf("Listening on %s", addr)
		err := http.ListenAndServe(addr, nil)
		catch(err)
	}()

	beltSize, ok := cfg.Get("belt.size").(int64)
	if !ok {
		log.Fatal("Missing belt.size in config.toml")
	}
	for ; beltSize > 0; beltSize-- {
		go belt.Loop()
	}

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt)

	log.Printf("Received %s; exiting", <-sigCh)
}

func catch(err error) {
	if err != nil {
		panic(err)
	}
}
