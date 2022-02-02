// Package goublu Args manages the os-passed-in args between goublu itself and Ublu.
package main

// Args holds the split of the os arguments between goublu args and Ublu args.
type Args struct {
	Goubluargs string
	Ubluargs   []string
}

// NewArgs instances a Args from the os args.
func NewArgs(osargs []string) (args *Args) {
	args = &Args{
		Goubluargs: args.getGoubluArgs(osargs),
		Ubluargs:   args.getUbluArgs(osargs),
	}
	return args
}

func (a *Args) getGoubluArgs(osargs []string) (goubluargs string) {
	goubluargs = ""
	if len(osargs) > 2 {
		if osargs[1] == "-g" {
			goubluargs = osargs[2]
		}
	}
	return goubluargs
}

func (a *Args) getUbluArgs(osargs []string) (ubluargs []string) {
	ubluargs = make([]string, 0, 0)
	if len(osargs) > 2 {
		if osargs[1] == "-g" {
			ubluargs = osargs[3:]
		} else {
			ubluargs = osargs[1:]
		}
	}
	return ubluargs
}
