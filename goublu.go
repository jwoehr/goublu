// GoUblu github.com/jwoehr/goublu
// goublu launches and serves as a better-than-Java console for
// https://github.com/jwoehr/ublu Ublu, a Java-coded domain-specific language
// for remote programming of IBM midrange and mainframe systems.
// Neither this project nor Ublu are associated with IBM.
package main

import (
	"bufio"
	"fmt"
	"github.com/jroimartin/gocui"
	"github.com/nsf/termbox-go"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	// "unicode/utf8"
)

var commandLineEditor gocui.Editor

var commandLineHistory = make([]string, 0, 20)

var commandPointer = -1

func appendHistory(l string) {
	commandLineHistory = append(commandLineHistory, strings.TrimSpace(l))
	commandPointer = len(commandLineHistory) - 1
}

func backHistory() string {
	var result string
	if commandPointer > -1 {
		result = commandLineHistory[commandPointer]
		commandPointer = max(commandPointer-1, -1)
	} else {
		commandPointer = -1 // so we can't over-decrement
	}
	return strings.TrimSpace(result)
}

func forwardHistory() string {
	var result string
	if -1 < commandPointer && commandPointer < len(commandLineHistory) {
		result = commandLineHistory[commandPointer]
		commandPointer = min(commandPointer+1, len(commandLineHistory)-1)
	} else {
		if commandPointer == -1 {
			commandPointer = 0
			result = commandLineHistory[commandPointer]
			commandPointer = min(commandPointer+1, len(commandLineHistory)-1)
		} else {
			commandPointer = len(commandLineHistory) - 1 // so we can't over-increment
		}
	}
	return strings.TrimSpace(result)
}

// How far from bottom we reserve our input area
const inputLineOffset = 2

// Obligatory layout redraw function
func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("ubluout", -1, -1, maxX, maxY-inputLineOffset); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Autoscroll = true
	}
	if v, err := g.SetView("ubluin", -1, maxY-inputLineOffset, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Autoscroll = true
		v.Editable = true
		v.Editor = commandLineEditor
		v.Wrap = true
		if _, err := g.SetCurrentView("ubluin"); err != nil {
			return err
		}
	}
	return nil
}

/*
// Exit via the gui instead of via Ublu
func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
*/

// Pipe input to Ublu
func ubluin(g *gocui.Gui, v *gocui.View, stdin io.WriteCloser) {
	var l string
	var err error
	cx, cy := v.Cursor()
	_, gy := g.Size()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}
	l = strings.TrimSpace(l)
	w, _ := g.View("ubluout")
	if l != "" {
		fmt.Fprint(w, "> "+l+"\n")
		io.WriteString(stdin, l+"\n")
	}
	appendHistory(l)
	v.Clear()
	v.MoveCursor(0-cx, (gy-inputLineOffset)-cy, false)
}

// Write to console output from Ublu
func ubluout(g *gocui.Gui, text string) {
	v, err := g.View("ubluout")
	if err != nil {
		// handle error
	}
	count := len(text)
	width, _ := g.Size()
	// This isn't right, we'll have to deal with rune width instead
	for i := 0; i < count; i = i + width {
		fmt.Fprint(v, text[i:min(count-1, i+width)])
		if i < count-1 {
			fmt.Fprint(v, "\n")
		}
	}
	termbox.Interrupt()
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func main() {

	// Prepare command
	myCmds := []string{"-jar", "/opt/ublu/ublu.jar", "-g", "--"}
	ubluArgs := append(myCmds, os.Args[1:]...)
	cmd := exec.Command("java", ubluArgs...)

	// Pipes
	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	defer stdout.Close()
	defer stderr.Close()

	// Readers
	outreader := bufio.NewReader(stdout)
	errreader := bufio.NewReader(stderr)

	// cogui
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	} else {
		g.Mouse = true
	}

	// Deliver Ublu's stdout
	go func() {
		for {
			text, _ := outreader.ReadString('\n')
			ubluout(g, text)
		}
	}()

	// Deliver Ublu's stderr
	go func() {
		for {
			text, _ := errreader.ReadString('\n')
			ubluout(g, text)
		}
	}()

	commandLineEditor = gocui.EditorFunc(func(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
		gx, gy := g.Size()
		cx, cy := v.Cursor()
		text, _ := v.Line(cy)

		// Shut up compiler
		gx = gx
		cy = cy

		switch {
		case ch != 0 && mod == 0:
			v.EditWrite(ch)
		case key == gocui.KeySpace:
			v.EditWrite(' ')
		case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
			v.EditDelete(true)
		case key == gocui.KeyDelete:
			v.EditDelete(false)
		case key == gocui.KeyInsert:
			v.Overwrite = !v.Overwrite
		case key == gocui.KeyEnter:
			ubluin(g, v, stdin)
		case key == gocui.KeyArrowDown:
			v.MoveCursor(0-cx, 0, false)
			v.Clear()
			for _, ch := range forwardHistory() {
				v.EditWrite(ch)
			}
		case key == gocui.KeyArrowUp:
			v.MoveCursor(0-cx, 0, false)
			v.Clear()
			for _, ch := range backHistory() {
				v.EditWrite(ch)
			}
		case key == gocui.KeyArrowLeft:
			v.MoveCursor(-1, 0, false)
		case key == gocui.KeyArrowRight:
			v.MoveCursor(1, 0, false)
		case key == gocui.KeyCtrlA:
			v.MoveCursor(0-cx, 0, false)
		case key == gocui.KeyCtrlB:
			v.MoveCursor(-1, 0, false)
		case key == gocui.KeyCtrlE:
			v.MoveCursor(len(text)-cx, 0, false)
		case key == gocui.KeyCtrlF:
			v.MoveCursor(1, 0, false)
		case key == gocui.KeyCtrlK:
			// this isn't quite correct but sorta works
			for i := cy; i < gy; i++ {
				v.EditDelete(false)
			}
		}
	})

	// defer g.Close()

	g.Cursor = true
	g.SetManagerFunc(layout)

	go func() {
		if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
			log.Panicln(err)
		}
	}()

	cmd.Run()

	g.Close()
	fmt.Println("Ublu has exited.")
	fmt.Println("Goodbye from Goublu!")
}
