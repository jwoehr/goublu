// History provides command line history in the goublu Ublu input area
package goublu

import "strings"

type History struct {
	commandLines   []string
	commandPointer int
}

func NewHistory() *History {
	h := &History{}
	h.commandLines = make([]string, 0, 20)
	h.commandPointer = -1
	return h
}

func (h *History) Append(l string) {
	h.commandLines = append(h.commandLines, strings.TrimSpace(l))
	h.commandPointer = len(h.commandLines) - 1
}

func (h *History) Back() string {
	var result string
	if h.commandPointer > -1 {
		result = h.commandLines[h.commandPointer]
		h.commandPointer = Max(h.commandPointer-1, -1)
	} else {
		h.commandPointer = -1 // so we can't over-decrement
	}
	return strings.TrimSpace(result)
}

func (h *History) Forward() string {
	var result string
	if -1 < h.commandPointer && h.commandPointer < len(h.commandLines) {
		result = h.commandLines[h.commandPointer]
		h.commandPointer = Min(h.commandPointer+1, len(h.commandLines)-1)
	} else {
		if h.commandPointer == -1 {
			h.commandPointer = 0
			if len(h.commandLines) > 0 {
				result = h.commandLines[h.commandPointer]
			}
			h.commandPointer = Min(h.commandPointer+1, len(h.commandLines)-1)
		} else {
			h.commandPointer = len(h.commandLines) - 1 // so we can't over-increment
		}
	}
	return strings.TrimSpace(result)
}

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
