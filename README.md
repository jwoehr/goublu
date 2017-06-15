# goublu
Go language front end to provide a better console interface to jwoehr/ublu

Works rudimentarily.

Usage
* Build:  go -build goublu.go
* Invoke: goublu ublu_arg ublu_arg ...

* Assumes Ublu is found in /opt/ublu/ublu.jar
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

Jack Woehr 2017-06-15
