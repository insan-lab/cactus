// Copyright 2014 The Cactus Authors. All rights reserved.

package sandbox

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"time"
)

// DockerImage is the Docker image used for sandboxed execution.
var DockerImage = "cactus-sandbox"

type dockerCell struct {
	dir string
}

// NewDockerCell creates a Cell backed by Docker containers.
func NewDockerCell() (Cell, error) {
	dir, err := ioutil.TempDir("", "cactus-docker-")
	if err != nil {
		return nil, err
	}
	return &dockerCell{dir: dir}, nil
}

func (d *dockerCell) Command(name string, args ...string) *Cmd {
	return NewCmd(&dockerRunner{
		dir:  d.dir,
		name: name,
		args: args,
	})
}

func (d *dockerCell) Create(name string) (*os.File, error) {
	return os.Create(path.Join(d.dir, name))
}

func (d *dockerCell) Dispose() error {
	fmt.Println("-- ", d.dir)
	return os.RemoveAll(d.dir)
}

type dockerRunner struct {
	dir    string
	name   string
	args   []string
	limits Resources
	usages Resources

	stdinW  *io.PipeWriter
	stdinR  *io.PipeReader
	stdoutW *io.PipeWriter
	stdoutR *io.PipeReader
	stderrW *io.PipeWriter
	stderrR *io.PipeReader
}

func (r *dockerRunner) SetLimits(res Resources) {
	r.limits = res
}

func (r *dockerRunner) Usages() Resources {
	return r.usages
}

func (r *dockerRunner) StdinPipe() (io.WriteCloser, error) {
	r.stdinR, r.stdinW = io.Pipe()
	return r.stdinW, nil
}

func (r *dockerRunner) StdoutPipe() (io.ReadCloser, error) {
	r.stdoutR, r.stdoutW = io.Pipe()
	return r.stdoutR, nil
}

func (r *dockerRunner) StderrPipe() (io.ReadCloser, error) {
	r.stderrR, r.stderrW = io.Pipe()
	return r.stderrR, nil
}

func (r *dockerRunner) Run() error {
	cidFile := path.Join(r.dir, ".cidfile")
	os.Remove(cidFile)

	dockerArgs := []string{
		"run", "-i",
		"--rm",
		"--network", "none",
		"--cidfile", cidFile,
		"-v", r.dir + ":/work",
		"-w", "/work",
	}

	if r.limits.Cpu > 0 {
		cpuSec := float64(r.limits.Cpu) / float64(time.Second)
		dockerArgs = append(dockerArgs, "--cpus", fmt.Sprintf("%.2f", cpuSec))
	}

	if r.limits.Memory > 0 {
		dockerArgs = append(dockerArgs, "--memory", fmt.Sprintf("%d", r.limits.Memory))
		dockerArgs = append(dockerArgs, "--memory-swap", fmt.Sprintf("%d", r.limits.Memory))
	}

	dockerArgs = append(dockerArgs, DockerImage, r.name)
	dockerArgs = append(dockerArgs, r.args...)

	fmt.Println(dockerArgs)

	cmd := exec.Command("docker", dockerArgs...)

	if r.stdinR != nil {
		cmd.Stdin = r.stdinR
	}
	if r.stdoutW != nil {
		cmd.Stdout = r.stdoutW
	}
	if r.stderrW != nil {
		cmd.Stderr = r.stderrW
	}

	start := time.Now()
	err := cmd.Start()
	if err != nil {
		r.closePipeWriters()
		return err
	}

	waitErr := cmd.Wait()
	wallTime := time.Since(start)

	r.closePipeWriters()

	r.usages.Cpu = wallTime
	r.collectUsages(cidFile)

	return waitErr
}

func (r *dockerRunner) closePipeWriters() {
	if r.stdoutW != nil {
		r.stdoutW.Close()
	}
	if r.stderrW != nil {
		r.stderrW.Close()
	}
}

func (r *dockerRunner) collectUsages(cidFile string) {
	cidBytes, err := ioutil.ReadFile(cidFile)
	if err != nil {
		return
	}
	cid := string(cidBytes)
	if cid == "" {
		return
	}

	out, err := exec.Command("docker", "inspect", cid).Output()
	if err != nil {
		return
	}

	var info []struct {
		State struct {
			OOMKilled bool `json:"OOMKilled"`
		} `json:"State"`
	}
	if err := json.Unmarshal(out, &info); err != nil || len(info) == 0 {
		return
	}

	if info[0].State.OOMKilled {
		r.usages.Memory = r.limits.Memory
	}
}
