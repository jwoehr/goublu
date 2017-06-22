// HelpReq is the F1 Help requestor
package goublu

import (
	"github.com/jwoehr/gocui"
)

const helpstring =
    "* Invoke: ./goublu [-g GoubluOpt1=SomeThing:GoubluOpt2=Other:...] ublu_arg ublu_arg ...\n" +
	"	* If the first argument to goublu is -g then the next element in the command line is assumed\n" +
	"	to be a string of Goublu property-like options of the form Opt=Value, each option separated from\n" +
	"	the next by : . All remaining commandline arguments are passed to Ublu. The Goublu options and their\n" +
	"	values are case-sensitive and are as follows: \n" +
	"		* UbluDir\n" +
	"			* abs path to dir where ublu.jar resides, default /opt/ublu\n" +
	"		* SaveOutDir\n" +
	"			* abs path to where pressing F4 saves the output text, default /tmp\n" +
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
	"		* Delete\n"+
	"* History\n"+
	"	* Up-arrow previous command\n"+
	"	* Down-arrow next command\n"+
	"	* PgUp first command\n"+
	"	* PgDn last command\n"+
	"	* F2 shows entire session's output\n"+
	"	* F4 saves the entire session's output to a file `/tmp/goublu.out.`_xxx..._\n"+
	"	* Output announces the save file name\n"+
	"		* You can do this as many times as you like during a session, a new file is created each time.\n"

func NewHelpReq(um *UbluManager, g *gocui.Gui) *ReqManager {
	return NewReqManager(um, g, 110, 42, " Goublu Help - F3 to exit ", gocui.ColorDefault, gocui.ColorDefault, helpstring)
}
