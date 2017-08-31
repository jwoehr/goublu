// Package goublu UbluManager is the gocui manager for Ublu input and output.
package goublu

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"

	"github.com/jwoehr/gocui"
	"github.com/jwoehr/termbox-go"
)

// How far from bottom we reserve our input area
const inputLineOffset = 3

// UbluManager is a gocui Manager for Ublu io.
type UbluManager struct {
	U                 *Ublu
	G                 *gocui.Gui
	Opts              *Options
	Hist              *History
	CommandLineEditor gocui.Editor
	Completor         *Completor
	ExitChan          chan string
	Dialoging         bool
}

// NewUbluManager instances a new manager.
func NewUbluManager(ublu *Ublu, g *gocui.Gui, opts *Options, hist *History) (um *UbluManager) {
	um = &UbluManager{
		U:         ublu,
		G:         g,
		Opts:      opts,
		Hist:      hist,
		Completor: NewCompletor(),
		ExitChan:  make(chan string),
		Dialoging: false,
	}
	um.CommandLineEditor = gocui.EditorFunc(func(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
		cx, cy := v.Cursor()
		text, _ := v.Line(cy)
		text = strings.Trim(strings.TrimSpace(text), "\000")

		switch {
		case ch != 0 && mod == 0:
			v.EditWrite(ch)
		case ch == 'b' && mod == gocui.ModAlt:
			backWord(v, cx, text)
		case ch == 'f' && mod == gocui.ModAlt:
			foreWord(v, cx, text)
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
			replaceLine(v, cx, um.Hist.Forward())
		case key == gocui.KeyArrowUp:
			replaceLine(v, cx, um.Hist.Back())
		case key == gocui.KeyPgup:
			replaceLine(v, cx, um.Hist.First())
		case key == gocui.KeyPgdn:
			replaceLine(v, cx, um.Hist.Last())
		case key == gocui.KeyArrowLeft:
			v.MoveCursor(-1, 0, false)
		case key == gocui.KeyArrowRight:
			v.MoveCursor(1, 0, false)
		case key == gocui.KeyCtrlSpace:
			replaceLine(v, cx, um.tryComplete(text))
		case key == gocui.KeyCtrlA || key == gocui.KeyHome:
			v.MoveCursor(0-cx, 0, false)
		case key == gocui.KeyCtrlB:
			v.MoveCursor(-1, 0, false)
		case key == gocui.KeyCtrlD:
			delWord(v, cx, text)
		case key == gocui.KeyCtrlE || key == gocui.KeyEnd:
			v.MoveCursor(len(text)-cx, 0, false)
		case key == gocui.KeyCtrlF:
			v.MoveCursor(1, 0, false)
		case key == gocui.KeyCtrlK:
			// this isn't quite correct but sorta works
			for i := cy; i < len(text); i++ {
				v.EditDelete(false)
			}
		case key == gocui.KeyF1:
			rm := NewHelpReq(um, um.G)
			rm.StartReq()
		case key == gocui.KeyF2:
			rm := NewAllOutReq(um, um.G)
			rm.StartReq()
		case key == gocui.KeyF3:
			um.doExitDialog(g, 9, 3)
		case key == gocui.KeyF4:
			f, err := ioutil.TempFile(um.Opts.SaveOutDir, "goublu.out.")
			if err != nil {
				log.Panicln(err)
			}
			um.Ubluout(um.G, "Saving output to "+f.Name()+"\n")
			f.Write([]byte(um.Hist.AllOut))
			f.Close()
		case key == gocui.KeyF5:
			replaceLine(v, cx, um.tryExpand(text))
		case key == gocui.KeyF9:
			replaceLine(v, cx, um.Hist.BackWrap())
		case key == gocui.MouseLeft:
		case key == gocui.MouseMiddle:
		case key == gocui.MouseRight:
		case key == gocui.MouseRelease:
		case key == gocui.MouseWheelUp:
		case key == gocui.MouseWheelDown:
		}
		if key != gocui.KeyCtrlSpace {
			um.Completor.Valid = false
		}
	})
	go func() {
		var exitMsg string
		for {
			exitMsg = <-um.ExitChan
			if exitMsg == "YES" {
				um.G.Close()
				fmt.Printf("Ublu exiting quickly on F3. Try to use 'bye' instead.\n")
				um.U.QuickExit(3)
			}
		}
	}()
	return um
}

// Ubluin pipes input to Ublu.
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

// Ubluout writes to console output from Ublu.
func (um *UbluManager) Ubluout(g *gocui.Gui, text string) {
	v, err := g.View("Ubluout")
	if err != nil {
		// handle error
	}
	fmt.Fprint(v, text)
	um.Hist.AppendAllOut(text)
	termbox.Interrupt()
}

// Layout is the obligatory gocui layout redraw function
func (um *UbluManager) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("Ubluout", 0, 0, maxX-1, maxY-inputLineOffset); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " Ublu Output  [F1 Help] [F2 Review] [F3 Quit?] [F4 Save] [F5 Macro] [F9 Prev] "
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
	if !um.Dialoging {
		if _, err := g.SetCurrentView("Ubluin"); err != nil {
			return err
		}
	}
	return nil
}

func (um *UbluManager) tryComplete(text string) (newtext string) {
	newtext = text
	if um.Completor.Valid == false {
		um.Completor.Clear()
	}
	if text != "" {
		words := strings.Fields(text)
		lastword := words[len(words)-1]
		candidate := um.Completor.Next()
		if candidate == "" {
			candidate = um.Completor.Complete(lastword)
		}
		if candidate != "" {
			newtext = text[0:strings.LastIndex(text, lastword)] + candidate
		}
	}
	return newtext
}

func (um *UbluManager) tryExpand(text string) (newtext string) {
	newtext = text
	if text != "" {
		words := strings.Fields(text)
		lastword := words[len(words)-1]
		candidate := um.Opts.Macros.Expand(lastword)
		if candidate != "" {
			newtext = text[0:strings.LastIndex(text, lastword)] + candidate
		}
	} else {
		um.Ubluout(um.G, "***\n")
		um.Ubluout(um.G, "Goublu Options:\n")
		um.Ubluout(um.G, um.Opts.AllOpts())
		um.Ubluout(um.G, "***\n")
		um.Ubluout(um.G, "Goublu Macros:\n")
		um.Ubluout(um.G, um.Opts.Macros.AllMacros())
		um.Ubluout(um.G, "***\n")
	}
	return newtext
}

func foreWord(v *gocui.View, cx int, text string) {
	nbfound := false
	if cx <= len(text) {
		i := cx
		for ; i < len(text); i++ {
			if text[i] != ' ' {
				nbfound = true
			}
			if nbfound && text[i] == ' ' {
				break
			}
		}
		v.MoveCursor(i-cx, 0, false)
	}
}

func backWord(v *gocui.View, cx int, text string) {
	nbfound := false
	if cx > 0 {
		i := Min(cx, len(text))
		for i = i - 1; i >= 0; i-- {
			if text[i] != ' ' {
				nbfound = true
			}
			if nbfound && text[i] == ' ' {
				break
			}
		}
		v.MoveCursor(i-cx, 0, false)
	}
}

func delWord(v *gocui.View, cx int, text string) {
	// if text[len(text)-1] != ' ' {
	// 	text = text + " "
	// }
	offset := Min(cx, len(text))
	head := text[:offset]
	tail := text[offset:]
	if len(tail) > 0 {
		tailSplit := strings.SplitN(tail, " ", 2)
		if len(tailSplit) > 1 {
			if len(tailSplit[1]) > 0 {
				if tailSplit[1][0] == ' ' {
					replaceLine(v, cx, head+tailSplit[1])
				} else {
					replaceLine(v, cx, head+" "+tailSplit[1])
				}
			}
		} else {
			replaceLine(v, cx, text[:Max(0, cx)])
		}
	} else {
		replaceLine(v, cx, text[:Max(0, cx)])
	}
}

func replaceLine(v *gocui.View, cx int, newtext string) {
	v.Clear()
	v.MoveCursor(0-cx, 0, false)
	for _, ch := range newtext {
		v.EditWrite(ch)
	}
}

func (um *UbluManager) doExitDialog(g *gocui.Gui, x int, y int) {
	NewTextLineDialog(um, g, x, y, "Exit", gocui.ColorRed, gocui.ColorBlack, "NO\nYES").StartDialog(um.ExitChan)
}
