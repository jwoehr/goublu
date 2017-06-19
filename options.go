// Options sorts out user options for Goublu.
package goublu

import (
	"bufio"
	"errors"
	"github.com/jroimartin/gocui"
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
func NewOptions() *Options {
	opts := &Options{
		UbluDir:    "/opt/ublu/",
		SaveOutDir: "/tmp",
		BgColorIn:  gocui.ColorDefault,
		FgColorIn:  gocui.ColorDefault,
		BgColorOut: gocui.ColorDefault,
		FgColorOut: gocui.ColorDefault,
	}
	return opts
}

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
