// Copyright 2014 The Cactus Authors. All rights reserved.

package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"

	"github.com/pelletier/go-toml"

	"github.com/FurqanSoftware/cactus/belt"
	"github.com/FurqanSoftware/cactus/sandbox"
)

var cfg *toml.Tree

func main() {
	runtime.GOMAXPROCS((runtime.NumCPU() + 1) / 2)

	_, err := os.Stat("config.toml")
	if os.IsNotExist(err) {
		log.Print("Creating config.toml from sample")

		f2, err := configSampleFS.Open("config-sample.toml")
		catch(err)

		f, err := os.Create("config.toml")
		catch(err)
		_, err = io.Copy(f, f2)
		catch(err)

		err = f2.Close()
		catch(err)
		err = f.Close()
		catch(err)
	}

	cfg, err = toml.LoadFile("config.toml")
	catch(err)
	log.Print("Loaded config.toml")

	go func() {
		addr, ok := cfg.Get("core.addr").(string)
		if !ok {
			log.Fatal("Missing core.addr in config.toml")
		}

		log.Printf("Listening on %s", addr)
		err := http.ListenAndServe(addr, nil)
		catch(err)
	}()

	sandboxMode := "jail"
	if mode, ok := cfg.Get("sandbox.mode").(string); ok {
		sandboxMode = mode
		switch mode {
		case "jail":
			belt.NewCell = func() (sandbox.Cell, error) {
				return sandbox.NewJailCell()
			}
		case "docker":
			belt.NewCell = func() (sandbox.Cell, error) {
				return sandbox.NewDockerCell()
			}
			if image, ok := cfg.Get("sandbox.docker_image").(string); ok {
				sandbox.DockerImage = image
			}
		default:
			log.Fatalf("Unknown sandbox.mode value %q in config.toml (expected \"jail\" or \"docker\")", mode)
		}
	}
	log.Printf("Sandbox mode: %s", sandboxMode)

	beltSize, ok := cfg.Get("belt.size").(int64)
	if !ok {
		log.Fatal("Missing belt.size in config.toml")
	}
	log.Printf("Starting %d judging worker(s)", beltSize)
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
