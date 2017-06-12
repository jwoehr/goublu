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
	"github.com/nsf/termbox-go"
	"github.com/jwoehr/goublu/console"
)

func main() {
	
	var edit_box console.EditBox
	
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
	/* scanner := bufio.NewScanner(os.Stdin) */

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
	/*
		go func() {
			for scanner.Scan() {
				io.WriteString(stdin, scanner.Text()+"\n")
			}
			if err := scanner.Err(); err != nil {
				fmt.Fprintln(os.Stderr, "reading standard input:", err)
			}
		}()
	*/
	go func() {
		err := termbox.Init()
		if err != nil {
			panic(err)
		}
		defer termbox.Close()
		termbox.SetInputMode(termbox.InputEsc)

		console.Redraw_all()
	mainloop:
		for {
			switch ev := termbox.PollEvent(); ev.Type {
			case termbox.EventKey:
				switch ev.Key {
				case termbox.KeyEsc:
					break mainloop
				case termbox.KeyArrowLeft, termbox.KeyCtrlB:
					edit_box.MoveCursorOneRuneBackward()
				case termbox.KeyArrowRight, termbox.KeyCtrlF:
					edit_box.MoveCursorOneRuneForward()
				case termbox.KeyBackspace, termbox.KeyBackspace2:
					edit_box.DeleteRuneBackward()
				case termbox.KeyDelete, termbox.KeyCtrlD:
					edit_box.DeleteRuneForward()
				case termbox.KeyTab:
					edit_box.InsertRune('\t')
				case termbox.KeySpace:
					edit_box.InsertRune(' ')
				case termbox.KeyCtrlK:
					edit_box.DeleteTheRestOfTheLine()
				case termbox.KeyHome, termbox.KeyCtrlA:
					edit_box.MoveCursorToBeginningOfTheLine()
				case termbox.KeyEnd, termbox.KeyCtrlE:
					edit_box.MoveCursorToEndOfTheLine()
				case termbox.KeyEnter:
					io.WriteString(stdin, string(edit_box.Text) + "\n")
					edit_box.Empty()
					edit_box.MoveCursorToBeginningOfTheLine()
					// edit_box.Draw()
				default:
					if ev.Ch != 0 {
						edit_box.InsertRune(ev.Ch)
					}
				}
			case termbox.EventError:
				panic(ev.Err)
			}
			console.Redraw_all()
		}
	}()
	
	cmd.Run()
}
