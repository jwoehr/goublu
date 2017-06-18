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
	"github.com/jwoehr/goublu"
	"github.com/nsf/termbox-go"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

var commandLineEditor gocui.Editor
var allOut string

// How far from bottom we reserve our input area
const inputLineOffset = 3

// Obligatory layout redraw function
func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("ubluout", 0, 0, maxX-1, maxY-inputLineOffset); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Ublu Output"
		v.Autoscroll = true
		v.Wrap = true
	}
	if v, err := g.SetView("ubluin", 0, maxY-inputLineOffset, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Ublu Input"
		v.Autoscroll = true
		v.Editable = true
		v.Editor = commandLineEditor
		v.Wrap = true
	}
	if _, err := g.SetCurrentView("ubluin"); err != nil {
		return err
	}
	return nil
}

// Pipe input to Ublu
func ubluin(g *gocui.Gui, v *gocui.View, stdin io.WriteCloser, history *goublu.History) {
	var l string
	var err error
	cx, cy := v.Cursor()
	_, gy := g.Size()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}
	l = strings.Trim(strings.TrimSpace(l), "\000")
	ubluout(g, l+"\n")
	io.WriteString(stdin, l+"\n")
	if l != "" {
		history.Append(l)
	}
	v.Clear()
	v.MoveCursor(0-cx, (gy-inputLineOffset)-cy, false)
}

// Write to console output from Ublu
func ubluout(g *gocui.Gui, text string) {
	v, err := g.View("ubluout")
	if err != nil {
		// handle error
	}
	fmt.Fprint(v, text)
	allOut = allOut + text
	termbox.Interrupt()
}

func main() {

	history := goublu.NewHistory()

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
			ubluin(g, v, stdin, history)
			termbox.Interrupt() // for good luck
		case key == gocui.KeyArrowDown:
			v.Clear()
			v.MoveCursor(0-cx, 0, false)
			for _, ch := range history.Forward() {
				v.EditWrite(ch)
			}
		case key == gocui.KeyArrowUp:
			v.Clear()
			v.MoveCursor(0-cx, 0, false)
			for _, ch := range history.Back() {
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
		case key == gocui.KeyF4:
			f, err := ioutil.TempFile("", "goublu.out.")
			if err != nil {
				log.Panicln(err)
			}
			ubluout(g, "Saving output to "+f.Name()+"\n")
			f.Write([]byte(allOut))			
			f.Close()
		}
	})

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
