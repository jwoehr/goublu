// Package goublu UbluManager is the gocui manager for Ublu input and output.
package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"

	"github.com/jroimartin/gocui"
	"github.com/nsf/termbox-go"
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
	CompileDate       string
	Version           string
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
	um.CommandLineEditor = um.createCommandLineEditor()
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

// createCommandLineEditor creates the editor function with panic recovery.
func (um *UbluManager) createCommandLineEditor() gocui.Editor {
	return gocui.EditorFunc(func(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Editor panic recovered: %v", r)
				um.Ubluout(um.G, fmt.Sprintf("Editor error: %v\n", r))
			}
		}()

		cx, cy := v.Cursor()
		text, _ := v.Line(cy)
		text = strings.Trim(strings.TrimSpace(text), "\000")

		// Handle character input
		if um.handleCharacterInput(v, ch, mod) {
			return
		}

		// Handle navigation keys
		if um.handleNavigationKeys(v, key, cx, text) {
			return
		}

		// Handle editing keys
		if um.handleEditingKeys(v, key, cx, cy, text) {
			return
		}

		// Handle function keys
		if um.handleFunctionKeys(v, key, cx, text) {
			return
		}

		// Handle mouse events (no-op for now)
		um.handleMouseEvents(key)

		// Invalidate completion if not Ctrl-Space
		if key != gocui.KeyCtrlSpace {
			um.Completor.Valid = false
		}
	})
}

// handleCharacterInput handles regular character input and Alt-modified characters.
func (um *UbluManager) handleCharacterInput(v *gocui.View, ch rune, mod gocui.Modifier) bool {
	cx, _ := v.Cursor()
	text, _ := v.Line(0)
	text = strings.Trim(strings.TrimSpace(text), "\000")

	switch {
	case ch != 0 && mod == 0:
		v.EditWrite(ch)
		return true
	case ch == 'b' && mod == gocui.ModAlt:
		backWord(v, cx, text)
		return true
	case ch == 'f' && mod == gocui.ModAlt:
		foreWord(v, cx, text)
		return true
	}
	return false
}

// handleNavigationKeys handles cursor movement keys.
func (um *UbluManager) handleNavigationKeys(v *gocui.View, key gocui.Key, cx int, text string) bool {
	switch key {
	case gocui.KeyArrowLeft:
		v.MoveCursor(-1, 0, false)
		return true
	case gocui.KeyArrowRight:
		v.MoveCursor(1, 0, false)
		return true
	case gocui.KeyArrowUp:
		replaceLine(v, cx, um.Hist.Back())
		return true
	case gocui.KeyArrowDown:
		replaceLine(v, cx, um.Hist.Forward())
		return true
	case gocui.KeyPgup:
		replaceLine(v, cx, um.Hist.First())
		return true
	case gocui.KeyPgdn:
		replaceLine(v, cx, um.Hist.Last())
		return true
	case gocui.KeyCtrlA, gocui.KeyHome:
		v.MoveCursor(0-cx, 0, false)
		return true
	case gocui.KeyCtrlE, gocui.KeyEnd:
		v.MoveCursor(len(text)-cx, 0, false)
		return true
	case gocui.KeyCtrlB:
		v.MoveCursor(-1, 0, false)
		return true
	case gocui.KeyCtrlF:
		v.MoveCursor(1, 0, false)
		return true
	}
	return false
}

// handleEditingKeys handles text editing operations.
func (um *UbluManager) handleEditingKeys(v *gocui.View, key gocui.Key, cx, cy int, text string) bool {
	switch key {
	case gocui.KeySpace:
		v.EditWrite(' ')
		return true
	case gocui.KeyBackspace, gocui.KeyBackspace2:
		v.EditDelete(true)
		return true
	case gocui.KeyDelete:
		v.EditDelete(false)
		return true
	case gocui.KeyInsert:
		v.Overwrite = !v.Overwrite
		return true
	case gocui.KeyEnter:
		um.Ubluin(um.G, v)
		termbox.Interrupt()
		return true
	case gocui.KeyCtrlD:
		delWord(v, cx, text)
		return true
	case gocui.KeyCtrlK:
		// Delete to end of line
		for i := cy; i < len(text); i++ {
			v.EditDelete(false)
		}
		return true
	case gocui.KeyCtrlSpace:
		replaceLine(v, cx, um.tryComplete(text))
		return true
	}
	return false
}

// handleFunctionKeys handles F1-F12 function keys.
func (um *UbluManager) handleFunctionKeys(v *gocui.View, key gocui.Key, cx int, text string) bool {
	switch key {
	case gocui.KeyF1:
		um.handleF1Help()
		return true
	case gocui.KeyF2:
		um.handleF2Review()
		return true
	case gocui.KeyF3:
		um.handleF3Exit()
		return true
	case gocui.KeyF4:
		um.handleF4Save()
		return true
	case gocui.KeyF5:
		replaceLine(v, cx, um.tryExpand(text))
		return true
	case gocui.KeyF9:
		replaceLine(v, cx, um.Hist.BackWrap())
		return true
	}
	return false
}

// handleF1Help shows help dialog.
func (um *UbluManager) handleF1Help() {
	rm := NewHelpReq(um, um.G)
	rm.StartReq()
}

// handleF2Review shows output review dialog.
func (um *UbluManager) handleF2Review() {
	rm := NewAllOutReq(um, um.G)
	rm.StartReq()
}

// handleF3Exit shows exit confirmation dialog.
func (um *UbluManager) handleF3Exit() {
	um.doExitDialog(um.G, 9, 3)
}

// handleF4Save saves output to a file with error handling.
func (um *UbluManager) handleF4Save() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Error saving output: %v", r)
			um.Ubluout(um.G, fmt.Sprintf("Error saving output: %v\n", r))
		}
	}()

	f, err := ioutil.TempFile(um.Opts.SaveOutDir, "goublu.out.")
	if err != nil {
		um.Ubluout(um.G, fmt.Sprintf("Error creating save file: %v\n", err))
		return
	}
	defer f.Close()

	um.Ubluout(um.G, "Saving output to "+f.Name()+"\n")
	if _, err := f.Write([]byte(um.Hist.AllOut)); err != nil {
		um.Ubluout(um.G, fmt.Sprintf("Error writing to file: %v\n", err))
		return
	}
}

// handleMouseEvents handles mouse input (currently no-op).
func (um *UbluManager) handleMouseEvents(key gocui.Key) bool {
	switch key {
	case gocui.MouseLeft, gocui.MouseMiddle, gocui.MouseRight,
		gocui.MouseRelease, gocui.MouseWheelUp, gocui.MouseWheelDown:
		return true
	}
	return false
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

// Ubluerr writes error output from Ublu with visual distinction.
func (um *UbluManager) Ubluerr(g *gocui.Gui, text string) {
	v, err := g.View("Ubluout")
	if err != nil {
		// handle error
	}
	// Write stderr with color distinction if supported
	fmt.Fprintf(v, "\033[31m%s\033[0m", text) // Red color for stderr
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
		um.Ubluout(um.G, "Goublu version "+um.Version+" compiled "+um.CompileDate+"\n")
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
