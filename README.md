# goublu
Go language front end to provide a better console interface to jwoehr/ublu

Works rudimentarily.

Usage
* Build via go -build goublu.go
* Invoke goublu arg arg ...

* Assumes Ublu is found in /opt/ublu/ublu.jar
* Basic line editing
	* Ctl-a move to head of line
	* Ctl-e move to end of line.
	* These work as you would expect:
		* Backspace
		* Left-arrow
		* Right-arrow
		* Insert
		* Delete

Jack Woehr 2017-06-15
