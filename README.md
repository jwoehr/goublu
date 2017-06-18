# goublu
Goublu is a [Go language](http://golang.org) front end that provides a better console interface to [Ublu](https://github.com/jwoehr/ublu) than the console support provided by Java.

Goublu is new and works with some limitations.

Report bugs or make feature requests in the [Issue Tracker](https://github.com/jwoehr/goublu/issues)

## Usage

* Fetch:  `go get github.com/jwoehr/goublu`
* Build:  `cd $GOPATH/src/github.com/jwoehr/goublu/main; go build goublu.go`
* Invoke: `./goublu ublu_arg ublu_arg ...`
* Assumes Ublu is found in `/opt/ublu/ublu.jar`
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
	* F4 saves the entire session's output to a file `/tmp/goublu.out.`_xxx..._
		* Output announces the save file name
		* You can do this as many times as you like during a session, a new file is created each time.
		
## Notes
	* The Ublu prompt appears on a line by itself in Goublu.
	* Goublu "history" is input line recall and is separate from Ublu's own `history` command.
	* Any Ublu application program output should include a newline as the Goublu output mechanism requires it.

## Bugs

* Serious
	* Ublu prompts for a password when an AS400 object is created with an invalid password and does not echo. However,
	Goublu **will indeed echo the password** even though Ublu's password prompt says the password will not be echoed.
* Trivial
	* None at this time

Jack Woehr 2017-06-18
