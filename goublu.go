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
	"io"
	"log"
	"os"
	"os/exec"
)

var DefaultEditor gocui.Editor

// Do layout in cogui loop.
// Most of the work actually in go functions in main()
func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("ubluout", -1, -1, maxX, maxY-4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Autoscroll = true
	}
	if v, err := g.SetView("ubluin", -1, maxY-4, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Editable = true
		v.Editor = DefaultEditor
		v.Wrap = true
		if _, err := g.SetCurrentView("ubluin"); err != nil {
			return err
		}
	}
	return nil
}

// Exit via the gui instead of via Ublu
func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func ubluout(g *gocui.Gui, text string) {
	v, err := g.View("ubluout")
	if err != nil {
		// handle error
	}
	fmt.Fprint(v, text)
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

	// DefaultEditor is the default editor.
	DefaultEditor = gocui.EditorFunc(func(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
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
			// v.EditNewLine()
			var l string
			var err error
			// cx, cy := v.Cursor()
			_, cy := v.Cursor()
			if l, err = v.Line(cy); err != nil {
				l = ""
			}
			w, _ := g.View("ubluout")
			fmt.Fprint(w, l)
			io.WriteString(stdin, l+"\n")
			v.Clear()
		case key == gocui.KeyArrowDown:
			v.MoveCursor(0, 1, false)
		case key == gocui.KeyArrowUp:
			v.MoveCursor(0, -1, false)
		case key == gocui.KeyArrowLeft:
			v.MoveCursor(-1, 0, false)
		case key == gocui.KeyArrowRight:
			v.MoveCursor(1, 0, false)
		}
	})

	defer g.Close()

	g.Cursor = true
	g.SetManagerFunc(layout)

	go func() {
		if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
			log.Panicln(err)
		}
	}()

	cmd.Run()
}
