// Package goublu Options parses and stores user options for Goublu.
package main

import (
	"bufio"
	"errors"
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/jroimartin/gocui"
)

// Options holds user options.
type Options struct {
	UbluDir    string
	SaveOutDir string
	BgColorIn  gocui.Attribute
	FgColorIn  gocui.Attribute
	BgColorOut gocui.Attribute
	FgColorOut gocui.Attribute
	JavaOpts   []string
	Macros     *MacroExpander
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
		JavaOpts:   make([]string, 0),
		Macros:     NewMacroExpander(),
	}
	return opts
}

// // AllOpts returns a string of all options.
// func (o *Options) AllOpts() (allopts string) {
// 	m := make(map[string]string)
// 	val := reflect.ValueOf(o).Elem()
//
// 	for i := 0; i < val.NumField(); i++ {
//
// 		valueField := val.Field(i)
// 		typeField := val.Type().Field(i)
//
// 		f := valueField.Interface()
// 		val := reflect.ValueOf(f)
// 		m[typeField.Name] = val.String()
// 	}
//
// 	for k, v := range m {
// 		if k != "Macros" {
// 			allopts += fmt.Sprintf("%s : %s\n", k, v)
// 		}
// 	}
// 	return
// }

// AllOpts returns a string of all options.
func (o *Options) AllOpts() (allopts string) {

	allopts += "JavaOpts : "
	for i := range o.JavaOpts {
		allopts += "\n\t" + o.JavaOpts[i]
	}
	allopts += "\n"
	allopts += "UbluDir : " + o.UbluDir + "\n"
	allopts += "SaveOutDir : " + o.SaveOutDir + "\n"
	allopts += "BgColorIn : " + NameFromColor(o.BgColorIn) + "\n"
	allopts += "FgColorIn : " + NameFromColor(o.FgColorIn) + "\n"
	allopts += "BgColorOut : " + NameFromColor(o.BgColorOut) + "\n"
	allopts += "FgColorOut : " + NameFromColor(o.FgColorOut) + "\n"
	return
}

// FromPropsFile parses and sets options from a property file.
func (o *Options) FromPropsFile(fname string) (err error) {
	f, err := os.Open(fname)
	if err == nil {
		propreader := bufio.NewReader(f)
		var prop string
		for {
			prop, err = propreader.ReadString('\n')
			if err != nil {
				break
			}
			o.FromPropString(prop)
		}
	}
	return err
}

// FromPropStrings parses and sets options from prop=val pairs separated by ':' .
// Note that prop string values in a props file can contain embedded : colons
// but prop string values on the command line may not because that's the prop
// seperator.
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

// FromPropString parses and sets an option from a single prop=val pair.
// Note that prop string values in a props file can contain embedded : colons
// but prop string values on the command line may not because that's the prop
// seperator.
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
			case "JavaOpt":
				o.JavaOpts = append(o.JavaOpts, val)
			case "Macro":
				o.Macros.AddFromProperty(val)
			default:
				err = errors.New("Unknown property")
			}
		} else {
			err = errors.New("Invalid property string")
		}
	}
	return err
}

// ColorFromName returns the named gocui color Attribute.
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

// NameFromColor returns the name of the gocui color Attribute.
func NameFromColor(color gocui.Attribute) (name string) {
	switch color {
	case gocui.ColorBlack:
		name = "ColorBlack"
	case gocui.ColorRed:
		name = "ColorRed"
	case gocui.ColorGreen:
		name = "ColorGreen"
	case gocui.ColorYellow:
		name = "ColorYellow"
	case gocui.ColorBlue:
		name = "ColorBlue"
	case gocui.ColorMagenta:
		name = "ColorMagenta"
	case gocui.ColorCyan:
		name = "ColorCyan"
	case gocui.ColorWhite:
		name = "ColorWhite"
	case gocui.ColorDefault:
		fallthrough
	default:
		name = "ColorDefault"
	}
	return
}
