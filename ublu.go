// Package goublu Ublu is the struct managing execution of the Java program Ublu.
package goublu

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

// Ublu is the controlling structure to run Java against ublu.jar with arguments.
type Ublu struct {
	Args      *Args
	Options   *Options
	Cmd       *exec.Cmd
	Stderr    io.ReadCloser
	Stdin     io.WriteCloser
	Stdout    io.ReadCloser
	ErrReader *bufio.Reader
	OutReader *bufio.Reader
}

// NewUblu returns an initialized *Ublu ready to run.
func NewUblu(args *Args, options *Options) (u *Ublu) {

	ubluFQP := filepath.Join(options.UbluDir, "ublu.jar")
	finfo, err := os.Stat(ubluFQP)

	if err != nil || finfo.IsDir() {
		fmt.Printf("%s does not exist or is not a file, please check your UbluDir property.\n", ubluFQP)
		return nil
	}

	myCmds := make([]string, 0)
	for o := range options.JavaOpts {
		myCmds = append(myCmds, options.JavaOpts[o])
	}
	myCmds = append(myCmds, "-jar")
	myCmds = append(myCmds, ubluFQP)
	myCmds = append(myCmds, "-g")
	myCmds = append(myCmds, "--")

	ubluArgs := append(myCmds, args.Ubluargs[:]...)

	cmd := exec.Command("java", ubluArgs...)

	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	outreader := bufio.NewReader(stdout)
	errreader := bufio.NewReader(stderr)

	u = &Ublu{
		Args:      args,
		Options:   options,
		Cmd:       cmd,
		Stderr:    stderr,
		Stdin:     stdin,
		Stdout:    stdout,
		OutReader: outreader,
		ErrReader: errreader,
	}

	return u
}

// Run runs the Ublu instance.
func (u *Ublu) Run() {
	u.Cmd.Run()
}

// Close does cleanup.
func (u *Ublu) Close() {
	u.Stdout.Close()
	u.Stderr.Close()
}

// QuickExit closes the pipes and exits with the given return code
func (u *Ublu) QuickExit(code int) {
	u.Close()
	os.Exit(code)
}
