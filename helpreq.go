// Package goublu HelpReq is the F1 Help requestor.
package goublu

import (
	"github.com/jwoehr/gocui"
)

const helpstring = "* Invoke: ./goublu [-g GoubluOpt1=SomeThing:GoubluOpt2=Other:...] ublu_arg ublu_arg ...\n" +
	"	* If the first argument to goublu is -g then the next element in the command line is assumed\n" +
	"	  to be a string of Goublu property-like options of the form Opt=Value, each option separated from\n" +
	"	  the next by : . All remaining commandline arguments are passed to Ublu. The Goublu options and their\n" +
	"	  values are read and interpreted in order, are case-sensitive and are as follows: \n" +
	"		* UbluDir\n" +
	"			* abs path to dir where ublu.jar resides, default /opt/ublu\n" +
	"		* JavaOpt\n" +
	"			* any option to the Java runtime, e.g, JavaOpt=-Dsomething=other (one option per JavaOpt line)\n" +
	"		* SaveOutDir\n" +
	"			* abs path to where pressing F4 saves the output text, default /tmp\n" +
	"		* PropsFile\n" +
	"			* abs path to a properties file containing these same option=value pairs\n" +
	"		* BgColorIn\n" +
	"			* Input background color, one of:\n" +
	"				* ColorBlack\n" +
	"				* ColorRed\n" +
	"				* ColorGreen\n" +
	"				* ColorYellow\n" +
	"				* ColorBlue\n" +
	"				* ColorMagenta\n" +
	"				* ColorCyan\n" +
	"				* ColorWhite\n" +
	"				* ColorDefault (default terminal colors)\n" +
	"		* FgColorIn\n" +
	"			* Input foreground color, as above\n" +
	"		* BgColorOut\n" +
	"			* Output background color, as above\n" +
	"		* FgColorOut\n" +
	"			* Output foreground color, as above\n" +
	"		* Macro=name freeform string of Ublu commands\n" +
	"			* Sets macro 'name' to 'freeform string of Ublu commands'\n" +
	"* Assumes in absence of property set as above that Ublu is found in /opt/ublu/ublu.jar\n" +
	"* Basic line editing\n" +
	"	* Ctl-a move to head of line\n" +
	"	* Ctl-b move one back.\n" +
	"	* Ctl-e move to end of line.\n" +
	"	* Ctl-f move one ahead.\n" +
	"	* Ctl-k delete to end of line.\n" +
	"		* This doesn't work entirely right if line is longer than view width.\n" +
	"	* These work as you would expect:\n" +
	"		* Backspace\n" +
	"		* Left-arrow\n" +
	"		* Right-arrow\n" +
	"		* Insert\n" +
	"		* Delete\n" +
	"* History\n" +
	"	* Up-arrow previous command\n" +
	"	* Down-arrow next command\n" +
	"	* PgUp first command\n" +
	"	* PgDn last command\n" +
	"* F1 shows Goublu help\n" +
	"* F2 shows entire session's output\n" +
	"* F3 offers a quick exit for when Ublu gets caught in a loop or network timeout\n" +
	"* F4 saves the entire session's output to a file `SaveOutDir/goublu.out.`_xxx..._\n" +
	"	* SaveOutDir set above as Goublu property, default is /tmp" +
	"	* Output announces the save file name\n" +
	"	* You can do this as many times as you like during a session, a new file is created each time.\n" +
	"* F5 expands last element on command line as macro you set in the properties file.\n" +
	"	* On empty line, F5 lists all macros.\n" +
	"* F9 rotates through previous commands wrapping.\n" +
	"* Ctrl-Space at the end of a partial command name rotates through completions, if any.\n"

// NewHelpReq returns a new help requester.
func NewHelpReq(um *UbluManager, g *gocui.Gui) *ReqManager {
	return NewReqManager(um, g, 110, 42, " Goublu Help - F3 to exit ", gocui.ColorDefault, gocui.ColorDefault, helpstring)
}
