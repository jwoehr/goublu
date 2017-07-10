// Completion does partial completion
package goublu

import (
	// "fmt"
	"strings"
)

// type Dashes []string

// type Completor has map[string] of Dash-Commands keyed by Ublu keywords and an
// array of the last set of candidate completions and the index of the last
// one tried
type Completor struct {
	CMap           map[string][]string
	LastCompletion []string
	NextIndex      int
}

// Ctor the Completor
func NewCompletor() (c *Completor) {
	c = &Completor{
		CMap:           make(map[string][]string),
		LastCompletion: make([]string, 0),
		NextIndex:      0,
	}
	c.CMap["as400"] = []string{"-to", "--"}
	c.CMap["ask"] = []string{}
	c.CMap["BREAK"] = []string{}
	c.CMap["bye"] = []string{}
	c.CMap["CALL"] = []string{}
	c.CMap["calljava"] = []string{}
	c.CMap["cim"] = []string{}
	c.CMap["cimi"] = []string{}
	c.CMap["collection"] = []string{}
	c.CMap["commandcall"] = []string{}
	c.CMap["const"] = []string{}
	c.CMap["cs"] = []string{}
	c.CMap["db"] = []string{}
	c.CMap["dbug"] = []string{}
	c.CMap["defun"] = []string{}
	c.CMap["desktop"] = []string{}
	c.CMap["dict"] = []string{}
	c.CMap["DO"] = []string{}
	c.CMap["dpoint"] = []string{}
	c.CMap["dq"] = []string{}
	c.CMap["dta"] = []string{}
	c.CMap["ELSE"] = []string{}
	c.CMap["eval"] = []string{}
	c.CMap["exit"] = []string{}
	c.CMap["file"] = []string{}
	c.CMap["FOR"] = []string{}
	c.CMap["oldftp"] = []string{}
	c.CMap["ftp"] = []string{}
	c.CMap["FUN"] = []string{}
	c.CMap["FUNC"] = []string{}
	c.CMap["gensh"] = []string{}
	c.CMap["help"] = []string{}
	c.CMap["histlog"] = []string{}
	c.CMap["h"] = []string{}
	c.CMap["host"] = []string{}
	c.CMap["history"] = []string{}
	c.CMap["IF"] = []string{}
	c.CMap["ifs"] = []string{}
	c.CMap["include"] = []string{}
	c.CMap["interpret"] = []string{}
	c.CMap["interpreter"] = []string{}
	c.CMap["jmx"] = []string{}
	c.CMap["job"] = []string{}
	c.CMap["joblist"] = []string{}
	c.CMap["joblog"] = []string{}
	c.CMap["jrnl"] = []string{}
	c.CMap["json"] = []string{}
	c.CMap["jvm"] = []string{}
	c.CMap["LOCAL"] = []string{}
	c.CMap["license"] = []string{}
	c.CMap["lifo"] = []string{}
	c.CMap["list"] = []string{}
	c.CMap["monitor"] = []string{}
	c.CMap["msg"] = []string{}
	c.CMap["msgq"] = []string{}
	c.CMap["num"] = []string{}
	c.CMap["objlist"] = []string{}
	c.CMap["objdesc"] = []string{}
	c.CMap["outq"] = []string{}
	c.CMap["ppl"] = []string{}
	c.CMap["printer"] = []string{}
	c.CMap["programcall"] = []string{}
	c.CMap["props"] = []string{}
	c.CMap["put"] = []string{}
	c.CMap["record"] = []string{}
	c.CMap["rs"] = []string{}
	c.CMap["savf"] = []string{}
	c.CMap["savef"] = []string{}
	c.CMap["savesys"] = []string{}
	c.CMap["server"] = []string{}
	c.CMap["sess"] = []string{}
	c.CMap["session"] = []string{}
	c.CMap["sleep"] = []string{}
	c.CMap["smapi"] = []string{}
	c.CMap["sock"] = []string{}
	c.CMap["spoolf"] = []string{}
	c.CMap["spoolflist"] = []string{}
	c.CMap["streamf"] = []string{}
	c.CMap["string"] = []string{}
	c.CMap["subsys"] = []string{}
	c.CMap["SWITCH"] = []string{}
	c.CMap["system"] = []string{}
	c.CMap["sysval"] = []string{}
	c.CMap["TASK"] = []string{}
	c.CMap["test"] = []string{}
	c.CMap["thread"] = []string{}
	c.CMap["THEN"] = []string{}
	c.CMap["THROW"] = []string{}
	c.CMap["tn5250"] = []string{}
	c.CMap["trace"] = []string{}
	c.CMap["TRY"] = []string{}
	c.CMap["tuple"] = []string{}
	c.CMap["usage"] = []string{}
	c.CMap["user"] = []string{}
	c.CMap["userlist"] = []string{}
	c.CMap["WHILE"] = []string{}
	c.CMap["!"] = []string{}
	c.CMap["#"] = []string{}
	c.CMap["#!"] = []string{}
	c.CMap["\\\\"] = []string{}
	return c
}

// Zero or more completion candidates returned
func (c *Completor) Complete(partial string) (candidate string) {
	candidates := make([]string, 0)
	// fmt.Printf("%d\n", len(c))
	for key, _ := range c.CMap {
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
	c.Set(candidates)
	if len(candidates) > 0 {
		candidate = c.Next()
	} else {
		candidate = ""
	}
	return candidate
}

// Set new list of candidate completions, set index to 0
func (c *Completor) Set(completion []string) {
	c.LastCompletion = completion
	c.NextIndex = 0
}

// Initialize candidate completions and index
func (c *Completor) Clear() {
	c.LastCompletion = make([]string, 0)
	c.NextIndex = 0
}

// Return next candidate completions and bump index
func (c *Completor) Next() (completion string) {
	if len(c.LastCompletion) > 0 {
		completion = c.LastCompletion[c.NextIndex]
		c.NextIndex = (c.NextIndex + 1) % len(c.LastCompletion)
	} else {
		completion = ""
	}
	return completion
}
