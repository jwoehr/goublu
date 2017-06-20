// Ublu is the struct managing execution of the Java program Ublu.
package goublu

import (
	"bufio"
	"io"
	"os/exec"
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

// Returns an initialized *Ublu ready to run.
func NewUblu(args *Args, options *Options) (u *Ublu) {

	myCmds := []string{"-jar", options.UbluDir + "/ublu.jar", "-g", "--"}
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

// Runs the Ublu instance.
func (u *Ublu) Run() {
	u.Cmd.Run()
}

// Does cleanup.
func (u *Ublu) Close() {
	u.Stdout.Close()
	u.Stderr.Close()
}
