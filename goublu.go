// GoUblu github.com/jwoehr/goublu
// goublu launches and serves as a better-than-Java console for
// https://github.com/jwoehr/ublu Ublu, a Java-coded domain-specific language
// for remote programming of IBM midrange and mainframe systems.
// Neither this project nor Ublu are associated with IBM.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jroimartin/gocui"
)

var compileDate string
var goubluVersion string

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "-v" || os.Args[1] == "--version") {
		fmt.Println("Goublu version", goubluVersion, "compiled", compileDate)
		os.Exit(0)
	}
	args := NewArgs(os.Args[:])
	options := NewOptions()
	options.FromPropStrings(args.Goubluargs)
	ublu := NewUblu(args, options)
	if ublu != nil {
		history := NewHistory()

		// cogui
		g, err := gocui.NewGui(gocui.OutputNormal)
		if err != nil {
			log.Panicln(err)
		} else {
			g.Mouse = true
		}

		um := NewUbluManager(ublu, g, options, history)
		g.SetManager(um)
		g.Cursor = true
		um.CompileDate = compileDate
		um.Version = goubluVersion

		// Deliver Ublu's stdout
		go func() {
			for {
				text, _ := ublu.OutReader.ReadString('\n')
				um.Ubluout(g, text)
			}
		}()

		// Deliver Ublu's stderr
		go func() {
			for {
				text, _ := ublu.ErrReader.ReadString('\n')
				um.Ubluout(g, text)
			}
		}()

		// Run the Gui
		go func() {
			if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
				log.Panicln(err)
			}
		}()

		ublu.Run()

		g.Close()
		ublu.Close()

		fmt.Println("Ublu has exited.")
	}

	fmt.Println("Goodbye from Goublu!")
}
