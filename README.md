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

## Building

* Fetch:  `go get -u github.com/jwoehr/goublu`
	* Or grab a [release of the source code](https://github.com/jwoehr/goublu/releases)
		* It's recommended you build from the latest release.
* Build:  `cd $GOPATH/src/github.com/jwoehr/goublu; ./make.sh`
  * Then copy `goublu` to your Golang binary path.
* Note: If you just do `go build` it won't compile version info into the image.
## Invoking
* Invoke: `./goublu [-v] [-g "GoubluOpt1=SomeThing:GoubluOpt2=Other:..."] ublu_arg ublu_arg ...`
  * If the first argument to goublu is `-v` then Goublu prints a version message and exits 0.
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

## Example setup
In `.bash_aliases` (or `.bashrc` or whatever):

```
alias gu='/home/jax/gopath/src/github.com/jwoehr/goublu/main/goublu -g PropsFile=/home/jax/.config/ublu/goublu.properties $*'
```

In the `PropsFile`:

```
#BgColorOut=ColorBlack
FgColorOut=ColorRed
UbluDir=/opt/ublu
SaveOutDir=.
JavaOpt=-Djavax.net.ssl.trustStore=/opt/ublu/keystore/ublutruststore
JavaOpt=-Dublu.includepath=/opt/ublu/examples:/opt/ublu/extensions
JavaOpt=-Dublu.usage.linelength=100
Macro=sys1 as400 -to @sys1 SYS1.FOO.COM myusrprf
Macro=in include
Macro=jl joblist -as400
Macro=db db -to @myDb -dbtype as400 -connect
Macro=ublutest /QSYS.LIB/UBLUTEST.LIB/
Macro=spfl spoolflist -as400
Macro=ul userlist -as400
Macro=ref desktop -browse file:///opt/ublu/userdoc/ubluref.html#
```

## Working in Goublu
* Basic line editing
	* Ctl-a move to head of line
	* Ctl-b move one back.
	* Ctl-e move to end of line.
	* Ctl-d delete to end of current word.
	* Ctl-f move one ahead.
	* Ctl-k delete to end of line.
		* This doesn't work entirely right if line is longer than view width.
	* Alt-b back a word.
	* Alt-f forward a word.
	* These work as you would expect:
		* Home
		* End
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
* F5 expands last element on command line as macro you set in the properties file.
	* On empty line, F5 lists Goublu version, compile date, and all Goublu options and macros.
* F9 rotates through previous commands wrapping.
* Ctrl-Space at the end of a partial command name rotates through completions, if any.

## Notes

* The Ublu prompt appears on a line by itself in Goublu.
* Goublu "history" is input line recall and is separate from Ublu's own `history` command.
* Any Ublu application program output should include a newline as the Goublu output mechanism requires it.
* This document as displayed on the project page always reflects the current state of the tree and may be in
advance of the release version.

## Bugs

* Serious
	* Ublu prompts for a password when an AS400 object is created with an invalid password and does not echo. However,
	Goublu **will indeed echo the password** even though Ublu's password prompt says the password will not be echoed.
* Trivial
	* Command lines longer that the view width of the input line behave erratically in response to edit commands.
	* On Mac OS X in Terminal, mouse actions fill the input line with escape sequences and do not otherwise work.

## The default branch has been renamed!

master is now named main

If you have a local clone, you can update it by running:
```
git branch -m master main
git fetch origin
git branch -u origin/main main
git remote set-head origin -a
```
Jack Woehr 2022-10-10
