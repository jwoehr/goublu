// ReqManager is the gocui manager for requestors.
package goublu

import (
	"fmt"
	"github.com/jwoehr/gocui"
	"strings"
)

// A gocui manager for Requestors.
type ReqManager struct {
	UM         *UbluManager
	G          *gocui.Gui
	X          int
	Y          int
	Title      string
	BgColor    gocui.Attribute
	FgColor    gocui.Attribute
	ReqEditor  gocui.Editor
	Text       string
	UbluinBuf  string
	UbluoutBuf string
}

// Instances a new manager.
func NewReqManager(um *UbluManager, g *gocui.Gui, x int, y int, title string, bgc gocui.Attribute, fgc gocui.Attribute, text string) (rm *ReqManager) {
	rm = &ReqManager{
		UM:      um,
		G:       g,
		X:       x,
		Y:       y,
		Title:   title,
		BgColor: bgc,
		FgColor: fgc,
		Text:    text,
	}
	rm.ReqEditor = gocui.EditorFunc(func(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
		gx, gy := rm.G.Size()
		cx, cy := v.Cursor()

		// Shut up compiler
		gx = gx
		gy = gy
		cx = cx
		cy = cy

		switch {
		case key == gocui.KeyArrowDown:
			v.MoveCursor(0, 1, false)
		case key == gocui.KeyArrowUp:
			v.MoveCursor(0, -1, false)
		case key == gocui.KeyArrowLeft:
			v.MoveCursor(-1, 0, false)
		case key == gocui.KeyArrowRight:
			v.MoveCursor(1, 0, false)
		case key == gocui.KeyPgup:
			ox, oy := v.Origin()
			v.SetOrigin(ox, oy-20)
		case key == gocui.KeyPgdn:
			ox, oy := v.Origin()
			v.SetOrigin(ox, oy+20)
		/*
			case key == gocui.MouseWheelUp:
				v.MoveCursor(0, -1, false)
				case key == gocui.MouseWheelDown:
				v.MoveCursor(0, 1, false)
		*/
		case key == gocui.KeyF3:
			rm.EndReq()
		}
	})

	return rm
}

func (rm *ReqManager) calcPos(maxX int, maxY int) (x1 int, y1 int, x2 int, y2 int) {
	x1 = Max((maxX-rm.X)/2, 0)
	x2 = maxX - x1
	y1 = Max((maxY-rm.Y)/2, 0)
	y2 = maxY - y1
	return x1, y1, x2, y2
}

// Unshows the requestor.
func (rm *ReqManager) EndReq() {
	rm.G.SetManager(rm.UM)
	rm.UM.Layout(rm.G)
	v, _ := rm.G.View("Ubluin")
	v.Clear()
	for _, ch := range strings.Trim(strings.TrimSpace(rm.UbluinBuf), "\000") {
		v.EditWrite(ch)
	}
	v, _ = rm.G.View("Ubluout")
	fmt.Fprintf(v, "%s\n", strings.Trim(strings.TrimSpace(rm.UbluoutBuf), "\000"))
}

// Shows the requestor.
func (rm *ReqManager) StartReq() {
	v, _ := rm.G.View("Ubluin")
	rm.UbluinBuf = v.Buffer()
	v, _ = rm.G.View("Ubluout")
	rm.UbluoutBuf = v.Buffer()
	rm.G.SetManager(rm.UM, rm)
}

// Obligatory layout redraw function
func (rm *ReqManager) Layout(g *gocui.Gui) error {
	x1, y1, x2, y2 := rm.calcPos(g.Size())
	if v, err := g.SetView(rm.Title, x1, y1, x2, y2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " " + rm.Title + " "
		v.Editable = true
		v.Editor = rm.ReqEditor
		v.Wrap = true
		v.BgColor = rm.BgColor
		v.FgColor = rm.FgColor
		fmt.Fprintf(v, "%s", rm.Text)
	}
	if _, err := g.SetCurrentView(rm.Title); err != nil {
		return err
	}
	return nil
}
