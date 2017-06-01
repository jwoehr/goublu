// Copyright *C* 2017 Jack J. Woehr
// All rights reserved.
// Use of this source code is governed by a BSD 3-clause license
// that can be found in the LICENSE file.
package main

import (
	"fmt"
	"io"
	"bufio"
	"os"
	"os/exec"
)

func main() {

	myCmds := []string {"-jar", "/opt/ublu/ublu.jar", "-g", "--"}
	ubluArgs := append(myCmds, os.Args[1:]...)
	cmd := exec.Command("java", ubluArgs...)
	stdin, _  := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	outreader := bufio.NewReader(stdout)
	errreader := bufio.NewReader(stderr)
/*	inreader  := bufio.NewReader(os.Stdin) */
	scanner   := bufio.NewScanner(os.Stdin)
	
	go func() {
		var text string
		for {
			text, _ = outreader.ReadString('\n')
			fmt.Print(text)
		}
		defer outreader.Close()
	}()
	
	go func() {
		var text string
		for {
			text, _ = errreader.ReadString('\n')
			fmt.Print(text)
		}
		defer errreader.Close()
	}()
/*	
	go func() {
		var text string
		for {
			text, _ = inreader.ReadString('\n')
			io.WriteString(stdin, text)  	  
		}
	}()
*/	
	go func() {
		for scanner.Scan() {
			io.WriteString(stdin, scanner.Text() + "\n")
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
		}
		defer scanner.Close()
	}()

	cmd.Run()
}