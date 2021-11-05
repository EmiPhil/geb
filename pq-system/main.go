// Decision Procedure for the pq-System as described in Gödel, Escher, Bach ch.2.
// Usage: pqs.exe ...strings where "...strings" is a list of possible theorems.
// Returns a table of the decisions made for each input string, including whether
// the input is an axiom or a theorem. Additionally, reports invalid inputs.

// The axiom schema of pq-System is xp-qx- whenever x is composed of hyphens only
// (and each x must stand for the same string).

// The only rule of production in pq-System is as follows: `Suppose x, y, and z
// all stand for particular strings containing only hyphens. And suppose that
// 'xpyqz' is known to be a theorem. Then 'xpy-qz-' is a theorem`.
package main

import (
	"fmt"
	"os"
)

// Position represents which part of the input string we are on: before the p, in
// between the p and the q, or after the q.
type Position int

const (
	LeftOfP Position = iota
	Center
	RightOfQ
)

// Lexer does all the heavy lifting. We read through a string and maintain our
// place here like in a grammatical lexer.
type Lexer struct {
	input    string
	place    int
	position Position
	valid    bool
	// leftOfP, betweenPAndQ, and rightOfQ represent the count of hyphens in each position
	leftOfP      int
	betweenPAndQ int
	rightOfQ     int
}

// Token is a shorthand for working with the possible runes in an easier way. On
// matching a rune to a character, we return the Token instead of the rune.
type Token int

const (
	TokenP Token = iota
	TokenQ
	TokenHyphen

	// TokenDone and TokenUnknown aren't valid runes but instead signify valuable
	// information to the lexer that we can't keep lexing.
	TokenDone
	TokenUnknown
)

// advance moves our place in the input string up by 1
func (l *Lexer) advance() {
	l.place++
}

// retreat moves our place in the input string down by 1
func (l *Lexer) retreat() {
	l.place--
}

// munch processes the **current** token and, in most cases, will advance our place by 1.
func (l *Lexer) munch() Token {
	defer l.advance()
	if l.place == len(l.input) {
		return TokenDone
	}

	r := rune(l.input[l.place])
	switch r {
	case 'p':
		fallthrough
	case 'P':
		return TokenP
	case 'q':
		fallthrough
	case 'Q':
		return TokenQ
	case '-':
		return TokenHyphen
	default:
		// Retreat our place: the earlier defer advance will increment our place even
		// though we didn't lex the current token.
		defer l.retreat()
		return TokenUnknown
	}
}

// process will munch through the entire input string, counting hyphens and
// updating the lexer's state as appropriate
func (l *Lexer) process() {
	token := l.munch()

	switch token {
	case TokenDone:
		l.valid = true
		return
	case TokenUnknown:
		l.valid = false
		return
	case TokenHyphen:
		switch l.position {
		case LeftOfP:
			l.leftOfP++
		case Center:
			l.betweenPAndQ++
		case RightOfQ:
			l.rightOfQ++
		}
	case TokenP:
		l.position = Center
	case TokenQ:
		l.position = RightOfQ
	}

	l.process()
}

// isAxiom verifies the axiom schema set out above by using the hyphen counts
func (l *Lexer) isAxiom() bool {
	return l.valid && l.leftOfP+1 == l.rightOfQ
}

// isTheorem verifies that this string is a theorem by using the hyphen counts
func (l *Lexer) isTheorem() bool {
	// all axioms are theorems
	if l.isAxiom() {
		return true
	}

	return l.valid && l.leftOfP+l.betweenPAndQ == l.rightOfQ
}

// NewLexer properly initializes a new Lexer from an input string
func NewLexer(s string) Lexer {
	return Lexer{input: s}
}

func main() {
	// input strings should be passed in as space separated strings
	strings := os.Args[1:]
	table := &Table{inputBorder: 7, headers: []string{"Input No.", "Valid", "Axiom", "Theorem", "Input"}}

	for n, s := range strings {
		lexer := NewLexer(s)
		lexer.process()

		// if the input string is very long, we extend the border of our table to match
		if len(lexer.input)+2 > table.inputBorder {
			table.inputBorder = len(lexer.input) + 2
		}

		// add the data of this input string to our table
		table.entries = append(table.entries, n+1, lexer.valid, lexer.isAxiom(), lexer.isTheorem(), lexer.input)
	}

	// pretty print the header, border, and data of the result table
	table.PrintBorder("┍", "┑", "━")
	table.PrintHeaders("│")
	table.PrintBorder("├", "┤", "─")
	table.PrintEntries("│")
	table.PrintBorder("└", "┘", "─")
}

// Table is used exclusively for formatting the results
type Table struct {
	inputBorder int
	headers     []string
	entries     []interface{}
}

// PrintBorder prints a solid line
func (t *Table) PrintBorder(leftCap, rightCap, borderCharacter string) {
	fmt.Printf(leftCap)

	for n, h := range t.headers {
		l := len(h) + 3
		if n == len(t.headers)-1 {
			l = len(h) + 2
		}

		for i := 0; i < l; i++ {
			fmt.Printf(borderCharacter)
		}
	}

	for i := 0; i < t.inputBorder-7; i++ {
		fmt.Printf(borderCharacter)
	}

	fmt.Printf("%s\n", rightCap)
}

// PrintHeaders prints the headers of the table using a given separator
func (t *Table) PrintHeaders(separator string) {
	for idx, h := range t.headers {
		fmt.Printf("%s %s ", separator, h)

		if idx == len(t.headers)-1 {
			for i := 0; i < t.inputBorder-7; i++ {
				fmt.Printf(" ")
			}
		}
	}
	fmt.Printf("%s\n", separator)
}

// PrintEntries prints each row of data, keeping in mind the need to start a new
// line every fifth cell
func (t *Table) PrintEntries(separator string) {
	for idx, cell := range t.entries {
		var format string
		switch idx % 5 {
		case 0:
			fmt.Printf(separator)
			format = " %9d %s"
		case 1:
			format = " %-5t %s"
		case 2:
			format = " %-5t %s"
		case 3:
			format = " %-7t %s"
		case 4:
			format = " %s "
			for i := len(cell.(string)) + 2; i < t.inputBorder; i++ {
				format += " "
			}
			format += "%s\n"

		}

		fmt.Printf(format, cell, separator)
	}
}
