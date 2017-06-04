// GoUblu github.com/jwoehr/goublu
// goublu launches ublu (https://github.com/jwoehr/ublu) a java-coded
// domain-specific language for remote programming of traditional IBM business
// systems. goublu means to supplement the weak console support of Java. 
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func main() {
	fmt.Println("GoUblu front end for Ublu")
	myCmds := []string{"-jar", "/opt/ublu/ublu.jar", "-g", "--"}
	ubluArgs := append(myCmds, os.Args[1:]...)
	cmd := exec.Command("java", ubluArgs...)
	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	outreader := bufio.NewReader(stdout)
	errreader := bufio.NewReader(stderr)
	/*	inreader  := bufio.NewReader(os.Stdin) */
	scanner := bufio.NewScanner(os.Stdin)

	go func() {
		var text string
		for {
			text, _ = outreader.ReadString('\n')
			fmt.Print(text)
		}
		defer stdout.Close()
	}()

	go func() {
		var text string
		for {
			text, _ = errreader.ReadString('\n')
			fmt.Print(text)
		}
		defer stderr.Close()
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
			io.WriteString(stdin, scanner.Text()+"\n")
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
		}
	}()

	cmd.Run()
}
