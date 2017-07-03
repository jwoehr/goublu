// Options parses and stores user options for Goublu.
package goublu

import (
	"bufio"
	"errors"
	/* debug */ // "fmt"
	"github.com/jwoehr/gocui"
	"os"
	"reflect"
	"regexp"
	"strings"
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
		UbluDir:    "/opt/ublu",
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
	if !strings.HasPrefix(prop, "#") {
		rx := regexp.MustCompile("=")
		s := rx.Split(prop, 2)
		if len(s) == 2 {
			val := strings.Trim(strings.TrimSpace(s[1]), "\012")
			/* debug */ // fmt.Println("Opt " + s[0] + " == " + val)
			switch s[0] {
			case "UbluDir":
				o.UbluDir = val
			case "SaveOutDir":
				o.SaveOutDir = val
			case "BgColorIn":
				o.BgColorIn = ColorFromName(val)
			case "FgColorIn":
				o.FgColorIn = ColorFromName(val)
			case "BgColorOut":
				o.BgColorOut = ColorFromName(val)
			case "FgColorOut":
				o.FgColorOut = ColorFromName(val)
			case "PropsFile":
				o.FromPropsFile(val)
			default:
				err = errors.New("Unknown property")
			}
		} else {
			err = errors.New("Invalid property string")
		}
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
