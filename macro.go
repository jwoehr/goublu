// Package goublu MacroExpander does expansion of macros.
package main

import "strings"

// MacroExpander is a map of macros
type MacroExpander struct {
	MMap map[string]string
}

// NewMacroExpander returns and empty expander
func NewMacroExpander() (m *MacroExpander) {
	m = &MacroExpander{MMap: make(map[string]string)}
	return m
}

// Add adds a macro and its expansion.
func (m *MacroExpander) Add(key string, expansion string) {
	m.MMap[key] = expansion
}

// AddFromProperty uses first word as key and the remainder as expansion.
func (m *MacroExpander) AddFromProperty(macroprop string) {
	if macroprop != "" {
		words := strings.Fields(macroprop)
		key := words[0]
		expansion := macroprop[len(key)+1:]
		m.Add(key, expansion)
	}
}

// Expand expands a macro or returns the key
func (m *MacroExpander) Expand(key string) (expansion string) {
	expansion = m.MMap[key]
	if expansion == "" {
		expansion = key
	}
	return
}

// AllMacros makes a displayable string of all macros and their expansions
func (m *MacroExpander) AllMacros() (allmacros string) {
	for key, value := range m.MMap {
		allmacros += key + " = " + value + "\n"
	}
	return
}
