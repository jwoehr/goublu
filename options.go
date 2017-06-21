// Options parses and stores user options for Goublu.
package goublu

import (
	"bufio"
	"errors"
	"github.com/jwoehr/gocui"
	"os"
	"reflect"
	"regexp"
)

// Options holds user options.
type Options struct {
	UbluDir    string
	SaveOutDir string
	BgColorIn  gocui.Attribute
	FgColorIn  gocui.Attribute
	BgColorOut gocui.Attribute
	FgColorOut gocui.Attribute
}

// NewOptions is a ctor with default options.
func NewOptions() (opts *Options) {
	opts = &Options{
		UbluDir:    "/opt/ublu/",
		SaveOutDir: "/tmp",
		BgColorIn:  gocui.ColorDefault,
		FgColorIn:  gocui.ColorDefault,
		BgColorOut: gocui.ColorDefault,
		FgColorOut: gocui.ColorDefault,
	}
	return opts
}

// Parses and sets options from a property file.
func (o *Options) FromPropsFile(fname string) (err error) {
	f, err := os.Open(fname)
	if err == nil {
		propreader := bufio.NewReader(f)

		for {
			prop, err := propreader.ReadString('\n')
			if err != nil {
				break
			}
			o.FromPropString(prop)
		}
	}
	return err
}

// Parses and sets options from prop=val pairs separated by ':' .
func (o *Options) FromPropStrings(props string) (err error) {
	err = nil
	rx := regexp.MustCompile(":")
	s := rx.Split(props, reflect.TypeOf(*o).NumField())
	for i := 0; i < len(s); i++ {
		err = o.FromPropString(s[i])
		if err != nil {
			break
		}
	}
	return err
}

// Parses and sets an option from a single prop=val pair.
func (o *Options) FromPropString(prop string) (err error) {
	err = nil
	rx := regexp.MustCompile("=")
	s := rx.Split(prop, 2)
	if len(s) == 2 {
		switch s[0] {
		case "UbluDir":
			o.UbluDir = s[1]
		case "SaveOutDir":
			o.SaveOutDir = s[1]
		case "BgColorIn":
			o.BgColorIn = ColorFromName(s[1])
		case "FgColorIn":
			o.FgColorIn = ColorFromName(s[1])
		case "BgColorOut":
			o.BgColorOut = ColorFromName(s[1])
		case "FgColorOut":
			o.FgColorOut = ColorFromName(s[1])
		default:
			err = errors.New("Unknown property")
		}
	} else {
		err = errors.New("Invalid property string")
	}
	return err
}

// Returns the named gocui color Attribute.
func ColorFromName(name string) (color gocui.Attribute) {
	switch name {
	case "ColorBlack":
		color = gocui.ColorBlack
	case "ColorRed":
		color = gocui.ColorRed
	case "ColorGreen":
		color = gocui.ColorGreen
	case "ColorYellow":
		color = gocui.ColorYellow
	case "ColorBlue":
		color = gocui.ColorBlue
	case "ColorMagenta":
		color = gocui.ColorMagenta
	case "ColorCyan":
		color = gocui.ColorCyan
	case "ColorWhite":
		color = gocui.ColorWhite
	case "ColorDefault":
		fallthrough
	default:
		color = gocui.ColorDefault
	}
	return color
}
