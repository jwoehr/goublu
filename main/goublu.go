// GoUblu github.com/jwoehr/goublu
// goublu launches and serves as a better-than-Java console for
// https://github.com/jwoehr/ublu Ublu, a Java-coded domain-specific language
// for remote programming of IBM midrange and mainframe systems.
// Neither this project nor Ublu are associated with IBM.
package main

import (
	"fmt"
	"github.com/jwoehr/gocui"
	"github.com/jwoehr/goublu"
	"log"
	"os"
)

func main() {
	args := goublu.NewArgs(os.Args[:])
	options := goublu.NewOptions()
	// options.FromPropStrings("BgColorOut=ColorRed:FgColorOut=ColorBlue")
	options.FromPropStrings(args.Goubluargs)
	ublu := goublu.NewUblu(args, options)
	history := goublu.NewHistory()

	// cogui
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	} else {
		g.Mouse = true
	}

	um := goublu.NewUbluManager(ublu, g, options, history)
	g.SetManager(um)
	g.Cursor = true

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
	fmt.Println("Goodbye from Goublu!")
}
