# goublu
Goublu is a [Go language](http://golang.org) front end that provides a better console interface to [Ublu](https://github.com/jwoehr/ublu) than the console support provided by Java.

![goublu_screenshot](https://user-images.githubusercontent.com/4604036/28322382-317d05fa-6b93-11e7-8457-b07eec2873af.png)

* Works on
	* OpenBSD
	* Linux
	* Mac OS X
* Works poorly on
	* Windows

Report bugs or make feature requests in the [Issue Tracker](https://github.com/jwoehr/goublu/issues)

## Usage

* Fetch:  `go get -u github.com/jwoehr/goublu`
* Build:  `cd $GOPATH/src/github.com/jwoehr/goublu/main; go build goublu.go`
* Invoke: `./goublu [-g "GoubluOpt1=SomeThing:GoubluOpt2=Other:..."] ublu_arg ublu_arg ...`
	* If the first argument to goublu is `-g` then the next element in the command line is assumed
	to be a string of Goublu property-like options of the form Opt=Value, each option separated from
	the next by `:` . All remaining commandline arguments are passed to Ublu. The Goublu options and their
	values are case-sensitive and are as follows:
		* `UbluDir`
			* abs path to dir where ublu.jar resides, default `/opt/ublu`
		* `JavaOpt`
			* any option to the Java runtime, e.g, `JavaOpt=-Dsomething=other` (one option per JavaOpt line)
		* `SaveOutDir`
			* abs path to where pressing F4 saves the output text, default `/tmp`
		* `PropsFile`
			* abs path to a properties file containing these same `option=value` pairs
		* `BgColorIn`
			* Input background color, one of:
				* `ColorBlack`
				* `ColorRed`
				* `ColorGreen`
				* `ColorYellow`
				* `ColorBlue`
				* `ColorMagenta`
				* `ColorCyan`
				* `ColorWhite`
				* `ColorDefault` (default terminal colors)
		* `FgColorIn`
			* Input foreground color, as above
		* `BgColorOut`
			* Output background color, as above
		* `FgColorOut`
			* Output foreground color, as above			
		* `Macro=name freeform string of Ublu commands`
			* Sets macro `name` to `freeform string of Ublu commands`
* Assumes in absence of property set as above that Ublu is found in `/opt/ublu/ublu.jar`
* Basic line editing
	* Ctl-a move to head of line
	* Ctl-b move one back.
	* Ctl-e move to end of line.
	* Ctl-f move one ahead.
	* Ctl-k delete to end of line.
		* This doesn't work entirely right if line is longer than view width.
	* These work as you would expect:
		* Backspace
		* Left-arrow
		* Right-arrow
		* Insert
		* Delete
* History
	* Up-arrow previous command
	* Down-arrow next command
	* PgUp first command
	* PgDn last command
* F1 shows Goublu help.
* F2 shows entire session's output
* F3 offers a quick exit for when Ublu gets caught in a loop or network timeout
* F4 saves the entire session's output to a file `SaveOutDir/goublu.out.`_xxx..._
	* SaveOutDir set above as Goublu property, default is /tmp
	* Output announces the save file name
	* You can do this as many times as you like during a session, a new file is created each time.
* F5 expands a macro you set in the properties file.
* F9 rotates through previous commands wrapping.
* Ctrl-Space at the end of a partial command name rotates through completions, if any.


## Notes

* The Ublu prompt appears on a line by itself in Goublu.
* Goublu "history" is input line recall and is separate from Ublu's own `history` command.
* Any Ublu application program output should include a newline as the Goublu output mechanism requires it.

## Bugs

* Serious
	* Ublu prompts for a password when an AS400 object is created with an invalid password and does not echo. However,
	Goublu **will indeed echo the password** even though Ublu's password prompt says the password will not be echoed.
* Trivial
	* Command lines longer that the view width of the input line behave erratically in response to edit commands.
	* On Mac OS X in Terminal, mouse actions fill the input line with escape sequences and do not otherwise work.

Jack Woehr 2017-07-18
