// Package goublu Dialog contains any user response dialogs for Goublu.
package goublu

import (
	"fmt"
	"strings"

	"github.com/jwoehr/gocui"
)

// TextLineDialog is what it sez.
type TextLineDialog struct {
	UM         *UbluManager
	G          *gocui.Gui
	X          int
	Y          int
	Title      string
	BgColor    gocui.Attribute
	FgColor    gocui.Attribute
	Text       string
	Result     string
	ResultChan chan string
}

// NewTextLineDialog instances a TextLineDialog.
func NewTextLineDialog(um *UbluManager, g *gocui.Gui, x int, y int, title string, bgc gocui.Attribute, fgc gocui.Attribute, text string) (tld *TextLineDialog) {
	tld = &TextLineDialog{
		UM:      um,
		G:       g,
		X:       x,
		Y:       y,
		Title:   title,
		BgColor: bgc,
		FgColor: fgc,
		Text:    text,
		Result:  "",
	}
	return
}

func (tld *TextLineDialog) calcPos(maxX int, maxY int) (x1 int, y1 int, x2 int, y2 int) {
	x1 = Max((maxX-tld.X)/2, 0)
	x2 = maxX - x1
	y1 = Max((maxY-tld.Y)/2, 0)
	y2 = maxY - y1
	return x1, y1, x2, y2
}

// Layout is the obligatory gocui layout redraw function.
func (tld *TextLineDialog) Layout(g *gocui.Gui) error {
	x1, y1, x2, y2 := tld.calcPos(g.Size())
	if v, err := g.SetView(tld.Title, x1, y1, x2, y2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " " + tld.Title + " "
		v.Wrap = false
		v.Highlight = true
		v.BgColor = tld.BgColor
		v.FgColor = tld.FgColor
		fmt.Fprintf(v, "%s", tld.Text)
	}
	if _, err := g.SetCurrentView(tld.Title); err != nil {
		return err
	}
	return nil
}

// StartDialog shows the dialog.
func (tld *TextLineDialog) StartDialog(exitChan chan string) error {
	tld.ResultChan = exitChan
	tld.G.SetManager(tld.UM, tld)
	tld.UM.Dialoging = true
	if err := tld.G.SetKeybinding(tld.Title, gocui.KeyEnter, gocui.ModNone, tld.EndDialog); err != nil {
		return err
	}
	if err := tld.G.SetKeybinding(tld.Title, gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := tld.G.SetKeybinding(tld.Title, gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}
	return nil
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	return nil
}

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy+1); err != nil {
			ox, oy := v.Origin()
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
	}
	return nil
}

// EndDialog unshows the dialog and stores the result.
func (tld *TextLineDialog) EndDialog(g *gocui.Gui, v *gocui.View) error {
	_, cy := v.Cursor()
	text, _ := v.Line(cy)
	tld.Result = strings.Trim(strings.TrimSpace(text), "\000")
	g.SetManager(tld.UM)
	tld.UM.Dialoging = false
	tld.UM.Layout(g)
	tld.ResultChan <- tld.Result
	return nil
}
