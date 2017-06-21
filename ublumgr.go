// UbluManager is the gocui manager for Ublu input and output.
package goublu

import (
	"fmt"
	"github.com/jwoehr/gocui"
	"github.com/jwoehr/termbox-go"
	"io"
	"io/ioutil"
	"log"
	"strings"
)

// How far from bottom we reserve our input area
const inputLineOffset = 3

// A gocui manager for Ublu io.
type UbluManager struct {
	U                 *Ublu
	G                 *gocui.Gui
	Opts              *Options
	Hist              *History
	CommandLineEditor gocui.Editor
}

// Instances a new manager.
func NewUbluManager(ublu *Ublu, g *gocui.Gui, opts *Options, hist *History) (um *UbluManager) {
	um = &UbluManager{
		U:    ublu,
		G:    g,
		Opts: opts,
		Hist: hist,
	}
	um.CommandLineEditor = gocui.EditorFunc(func(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
		gx, gy := um.G.Size()
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
			um.Ubluin(um.G, v)
			termbox.Interrupt() // for good luck
		case key == gocui.KeyArrowDown:
			v.Clear()
			v.MoveCursor(0-cx, 0, false)
			for _, ch := range um.Hist.Forward() {
				v.EditWrite(ch)
			}
		case key == gocui.KeyArrowUp:
			v.Clear()
			v.MoveCursor(0-cx, 0, false)
			for _, ch := range um.Hist.Back() {
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
		case key == gocui.KeyF1:
			rm := NewHelpReq(um, um.G)
			rm.StartReq()
		case key == gocui.KeyF2:
			rm := NewAllOutReq(um, um.G)
			rm.StartReq()
		case key == gocui.KeyF4:
			f, err := ioutil.TempFile(um.Opts.SaveOutDir, "goublu.out.")
			if err != nil {
				log.Panicln(err)
			}
			um.Ubluout(um.G, "Saving output to "+f.Name()+"\n")
			f.Write([]byte(um.Hist.AllOut))
			f.Close()
		}
	})

	return um
}

// Pipes input to Ublu
func (um *UbluManager) Ubluin(g *gocui.Gui, v *gocui.View) {
	var l string
	var err error
	cx, cy := v.Cursor()
	_, gy := g.Size()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}
	l = strings.Trim(strings.TrimSpace(l), "\000")
	um.Ubluout(g, l+"\n")
	io.WriteString(um.U.Stdin, l+"\n")
	if l != "" {
		um.Hist.Append(l)
	}
	v.Clear()
	v.MoveCursor(0-cx, (gy-inputLineOffset)-cy, false)
}

// Writes to console output from Ublu
func (um *UbluManager) Ubluout(g *gocui.Gui, text string) {
	v, err := g.View("Ubluout")
	if err != nil {
		// handle error
	}
	fmt.Fprint(v, text)
	um.Hist.AppendAllOut(text)
	termbox.Interrupt()
}

// Obligatory layout redraw function
func (um *UbluManager) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("Ubluout", 0, 0, maxX-1, maxY-inputLineOffset); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " Ublu Output  [F1 Goublu Help] [F2 Review Output] [F4 Save Output] "
		v.Autoscroll = true
		v.Wrap = true
		v.BgColor = um.Opts.BgColorOut
		v.FgColor = um.Opts.FgColorOut
	}
	if v, err := g.SetView("Ubluin", 0, maxY-inputLineOffset, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " Ublu Input "
		v.Autoscroll = true
		v.Editable = true
		v.Editor = um.CommandLineEditor
		v.Wrap = true
		v.BgColor = um.Opts.BgColorIn
		v.FgColor = um.Opts.FgColorIn
	}
	if _, err := g.SetCurrentView("Ubluin"); err != nil {
		return err
	}
	return nil
}
