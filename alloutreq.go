// Package goublu AllOutReq is the F2 Show All Output requestor
package main

import (
	"github.com/jroimartin/gocui"
)

// NewAllOutReq pops a view of all output in a Requestor
func NewAllOutReq(um *UbluManager, g *gocui.Gui) *ReqManager {
	return NewReqManager(um, g, 110, 42, " Ublu Output Review - F3 to exit ", um.Opts.BgColorOut, um.Opts.FgColorOut, um.Hist.AllOut)
}
