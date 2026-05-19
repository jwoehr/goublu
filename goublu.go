// GoUblu github.com/jwoehr/goublu
// goublu launches and serves as a better-than-Java console for
// https://github.com/jwoehr/ublu Ublu, a Java-coded domain-specific language
// for remote programming of IBM midrange and mainframe systems.
// Neither this project nor Ublu are associated with IBM.
package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/jroimartin/gocui"
)

var CompileDate string
var GoubluVersion string

// startStreamReader starts a goroutine to read from a stream and deliver output via handler.
// It respects context cancellation for graceful shutdown.
func startStreamReader(ctx context.Context, reader *bufio.Reader, streamName string, handler func(string)) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("%s reader panic: %v", streamName, r)
			}
		}()
		for {
			select {
			case <-ctx.Done():
				log.Printf("%s reader stopped by context cancellation", streamName)
				return
			default:
				text, err := reader.ReadString('\n')
				if err != nil {
					if err != io.EOF {
						log.Printf("Error reading %s: %v", streamName, err)
					}
					return
				}
				handler(text)
			}
		}
	}()
}

// initializeGUI creates and configures the GUI and UbluManager.
// Returns the GUI, UbluManager, and any error encountered.
func initializeGUI(ublu *Ublu, options *Options, history *History) (*gocui.Gui, *UbluManager, error) {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return nil, nil, err
	}

	g.Mouse = true
	g.Cursor = true

	um := NewUbluManager(ublu, g, options, history)
	g.SetManager(um)
	um.CompileDate = CompileDate
	um.Version = GoubluVersion

	return g, um, nil
}

// runApplication manages the application lifecycle including GUI and Ublu execution.
// Returns an error if either the GUI or Ublu encounters a fatal error.
func runApplication(ctx context.Context, g *gocui.Gui, um *UbluManager, ublu *Ublu) error {
	// Start stream readers with context cancellation support
	startStreamReader(ctx, ublu.OutReader, "stdout", func(text string) {
		um.Ubluout(g, text)
	})
	startStreamReader(ctx, ublu.ErrReader, "stderr", func(text string) {
		um.Ubluerr(g, text)
	})

	// Run the GUI in a goroutine and capture errors
	guiDone := make(chan error, 1)
	go func() {
		if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
			guiDone <- err
		}
		close(guiDone)
	}()

	// Run Ublu (blocks until exit)
	ubluErr := ublu.Run()

	// Wait for GUI to finish
	guiErr := <-guiDone

	// Cleanup
	g.Close()
	ublu.Close()

	// Return the first error encountered
	if guiErr != nil {
		return guiErr
	}
	return ubluErr
}

func main() {
	// Handle version flag with early return
	if len(os.Args) > 1 && (os.Args[1] == "-v" || os.Args[1] == "--version") {
		fmt.Println("Goublu version", GoubluVersion, "compiled", CompileDate)
		return
	}

	args := NewArgs(os.Args[:])
	options := NewOptions()
	options.FromPropStrings(args.Goubluargs)

	ublu := NewUblu(args, options)
	if ublu == nil {
		fmt.Println("Failed to initialize Ublu")
		fmt.Println("Goodbye from Goublu!")
		os.Exit(1)
	}

	history := NewHistory()

	// Initialize GUI
	g, um, err := initializeGUI(ublu, options, history)
	if err != nil {
		log.Printf("Failed to initialize GUI: %v", err)
		fmt.Println("Goodbye from Goublu!")
		os.Exit(2)
	}

	// Create context for graceful shutdown of stream readers
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Run the application
	if err := runApplication(ctx, g, um, ublu); err != nil {
		log.Printf("Application error: %v", err)
		fmt.Println("Goodbye from Goublu!")
		os.Exit(3)
	}

	fmt.Println("Ublu has exited.")
	fmt.Println("Goodbye from Goublu!")
}
