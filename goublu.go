package main

import (
	"fmt"
	"io"
	"bufio"
	"os"
	"os/exec"
)

func main() {
//	cmd := exec.Command("java", "-version")

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
	}()
	
	go func() {
		var text string
		for {
			text, _ = errreader.ReadString('\n')
			fmt.Print(text)
		}
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
	}()

	cmd.Run()
}