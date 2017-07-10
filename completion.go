// Completion does partial completion
package goublu

import (
	// "fmt"
	"strings"
)

// type Dashes []string

// type Completion map[string] of Dash-Commands keyed by Ublu keywords
type Completion map[string][]string

// Return the Completion map
func NewCompletion() (c Completion) {
	c = make(Completion)
	c["as400"] = []string{"-to", "--"}
	c["ask"] = []string{}
	c["BREAK"] = []string{}
	c["bye"] = []string{}
	c["CALL"] = []string{}
	c["calljava"] = []string{}
	c["cim"] = []string{}
	c["cimi"] = []string{}
	c["collection"] = []string{}
	c["commandcall"] = []string{}
	c["const"] = []string{}
	c["cs"] = []string{}
	c["db"] = []string{}
	c["dbug"] = []string{}
	c["defun"] = []string{}
	c["desktop"] = []string{}
	c["dict"] = []string{}
	c["DO"] = []string{}
	c["dpoint"] = []string{}
	c["dq"] = []string{}
	c["dta"] = []string{}
	c["ELSE"] = []string{}
	c["eval"] = []string{}
	c["exit"] = []string{}
	c["file"] = []string{}
	c["FOR"] = []string{}
	c["oldftp"] = []string{}
	c["ftp"] = []string{}
	c["FUN"] = []string{}
	c["FUNC"] = []string{}
	c["gensh"] = []string{}
	c["help"] = []string{}
	c["histlog"] = []string{}
	c["h"] = []string{}
	c["host"] = []string{}
	c["history"] = []string{}
	c["IF"] = []string{}
	c["ifs"] = []string{}
	c["include"] = []string{}
	c["interpret"] = []string{}
	c["interpreter"] = []string{}
	c["jmx"] = []string{}
	c["job"] = []string{}
	c["joblist"] = []string{}
	c["joblog"] = []string{}
	c["jrnl"] = []string{}
	c["json"] = []string{}
	c["jvm"] = []string{}
	c["LOCAL"] = []string{}
	c["license"] = []string{}
	c["lifo"] = []string{}
	c["list"] = []string{}
	c["monitor"] = []string{}
	c["msg"] = []string{}
	c["msgq"] = []string{}
	c["num"] = []string{}
	c["objlist"] = []string{}
	c["objdesc"] = []string{}
	c["outq"] = []string{}
	c["ppl"] = []string{}
	c["printer"] = []string{}
	c["programcall"] = []string{}
	c["props"] = []string{}
	c["put"] = []string{}
	c["record"] = []string{}
	c["rs"] = []string{}
	c["savf"] = []string{}
	c["savef"] = []string{}
	c["savesys"] = []string{}
	c["server"] = []string{}
	c["sess"] = []string{}
	c["session"] = []string{}
	c["sleep"] = []string{}
	c["smapi"] = []string{}
	c["sock"] = []string{}
	c["spoolf"] = []string{}
	c["spoolflist"] = []string{}
	c["streamf"] = []string{}
	c["string"] = []string{}
	c["subsys"] = []string{}
	c["SWITCH"] = []string{}
	c["system"] = []string{}
	c["sysval"] = []string{}
	c["TASK"] = []string{}
	c["test"] = []string{}
	c["thread"] = []string{}
	c["THEN"] = []string{}
	c["THROW"] = []string{}
	c["tn5250"] = []string{}
	c["trace"] = []string{}
	c["TRY"] = []string{}
	c["tuple"] = []string{}
	c["usage"] = []string{}
	c["user"] = []string{}
	c["userlist"] = []string{}
	c["WHILE"] = []string{}
	c["!"] = []string{}
	c["#"] = []string{}
	c["#!"] = []string{}
	c["\\\\"] = []string{}
	return c
}

// Zero or more completion candidates returned
func (c Completion) Complete(partial string) (candidates []string) {
	candidates = make([]string, 0)
	// fmt.Printf("%d\n", len(c))
	for key, _ := range c {
		// fmt.Printf("Partial is |%s|\n", partial)
		// fmt.Printf("Key is %s\n", key)
		if strings.HasPrefix(key, partial) {
			// fmt.Printf("Had prefix\n")
			candidates = append(candidates, key)
		} else {
			// fmt.Printf("NO prefix\n")
		}
	}
	// fmt.Printf("%d\n", len(candidates))
	return candidates
}
