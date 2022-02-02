// Package goublu History provides command line history in the goublu Ublu input area.
package main

// History keeps history lines from the goublu input line and output from the console out.
type History struct {
	commandLines   []string
	commandPointer int
	AllOut         string
}

// NewHistory is Ctor/0.
func NewHistory() *History {
	h := &History{}
	h.commandLines = make([]string, 0, 0)
	h.commandPointer = -1
	return h
}

// Append postpends a history line to the collection.
func (h *History) Append(l string) {
	h.commandLines = append(h.commandLines, l)
	h.commandPointer = len(h.commandLines) - 1
}

// Back returns the next backwards history line, decrementing the pointer for
// next history line request. The pointer is floored and does not wrap.
func (h *History) Back() string {
	var result string
	h.commandPointer = Min(len(h.commandLines)-1, Max(-1, h.commandPointer))
	if h.commandPointer >= 0 {
		result = h.commandLines[h.commandPointer]
		h.commandPointer = h.commandPointer - 1
	}
	return result
}

// BackWrap returns the next backwards history line, decrementing the pointer for
// next history line request. The pointer wraps.
func (h *History) BackWrap() string {
	var result string
	if len(h.commandLines) > 0 {
		if h.commandPointer > -1 {
			result = h.commandLines[h.commandPointer]
			h.commandPointer = (h.commandPointer - 1) % len(h.commandLines)
		} else {
			h.commandPointer = len(h.commandLines) - 1 // so we can't over-decrement
			if h.commandPointer > 0 {
				result = h.commandLines[h.commandPointer]
			}
			h.commandPointer = (h.commandPointer - 1) % len(h.commandLines)
		}
	}
	return result
}

// Forward returns the next forwards history line, incrementing the pointer for
// next history line request. The pointer is roofed and does not wrap.
func (h *History) Forward() string {
	var result string
	h.commandPointer = h.commandPointer + 1
	if h.commandPointer >= len(h.commandLines) {
		h.commandPointer = len(h.commandLines) - 1
	} else if h.commandPointer > -1 {
		result = h.commandLines[h.commandPointer]
	} else {
		h.commandPointer = -1
	}
	return result
}

// First returns the first history line if any, resetting the pointer if any.
func (h *History) First() string {
	var result string
	if len(h.commandLines) > 0 {
		h.commandPointer = 0
		result = h.commandLines[h.commandPointer]
	}
	return result
}

// Last returns the last history line if any, resetting the pointer if any.
func (h *History) Last() string {
	var result string
	if len(h.commandLines) > 0 {
		h.commandPointer = len(h.commandLines) - 1
		result = h.commandLines[h.commandPointer]
	}
	return result
}

// AppendAllOut adds more output text to the complete console out text.
func (h *History) AppendAllOut(text string) {
	h.AllOut = h.AllOut + text
}

// Max returns the max of two ints.
func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

// Min returns the min of two ints.
func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
