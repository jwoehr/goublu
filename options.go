// Options sorts out user options for Goublu.
package goublu

import (
	"github.com/jroimartin/gocui"
)

// Options holds user options.
type Options struct {
	Ubludir    string
	SaveOutDir string
	BgColorIn  gocui.Attribute
	FgColorIn  gocui.Attribute
	BgColorOut gocui.Attribute
	FgColorOut gocui.Attribute
}

// NewOptions is a ctor with default options.
func NewOptions() *Options {
	opts := &Options{
		Ubludir:    "/opt/ublu/",
		SaveOutDir: "/tmp",
		BgColorIn:  gocui.ColorDefault,
		FgColorIn:  gocui.ColorDefault,
		BgColorOut: gocui.ColorDefault,
		FgColorOut: gocui.ColorDefault,
	}
	return opts
}
