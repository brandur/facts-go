package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/scanner"

	"github.com/kardianos/osext"
)

type Fact struct {
	Front string
	Back  string
	Tags  []string
}

func dataFiles() ([]string, error) {
	filename, err := osext.Executable()
	if err != nil {
		return nil, err
	}

	path, err := filepath.Abs(filename + "/../categories/*.md")
	if err != nil {
		return nil, err
	}

	return filepath.Glob(path)
}

func parseMarkdown(str string) (string, []*Fact, error) {
	var s scanner.Scanner
	var tok rune

	// result fields
	var facts []*Fact
	var header string

	// current parsing state
	var buf bytes.Buffer
	var fact *Fact
	inHeader := false
	inNote := false
	indentCount := 0
	whitespaceCount := 0

	enableDebug := false
	if os.Getenv("FACTS_DEBUG") == "true" {
		enableDebug = true
	}

	// just makes output a little more succinct
	toB := func(b bool) string {
		if b {
			return "t"
		} else {
			return "f"
		}
	}

	debug := func(format string, a ...interface{}) {
		if enableDebug {
			a = append(a, indentCount, toB(inHeader), toB(inNote))
			fmt.Printf("scanner: "+format+"\t\t(indent=%v header=%v note=%v)\n", a...)
		}
	}

	debugPlain := func(format string, a ...interface{}) {
		if enableDebug {
			fmt.Printf("scanner: "+format+"\n", a...)
		}
	}

	s.Init(strings.NewReader(str))

	// set mode to 0 to avoid any attempt at string reading, which may cause
	// problems when we run into apostrophes
	s.Mode = scanner.ScanFloats | scanner.ScanIdents | scanner.ScanInts
	s.Whitespace = 1 << '\r'

	for tok != scanner.EOF {
		tok = s.Scan()

		switch s.TokenText() {
		case " ":
			debug("[SPACE]")

			if inHeader || inNote {
				buf.WriteString(s.TokenText())
			} else {
				whitespaceCount += 1
				if whitespaceCount == 4 {
					indentCount += 1
					whitespaceCount = 0
				}
			}

		// reset everything and handle a completed string
		case "\n":
			debug("\\n")

			// string is likely left with leading and possibly trailing
			// whitespace
			trimmed := strings.TrimSpace(buf.String())

			if inNote && trimmed != "" {
				if indentCount == 0 {
					debugPlain("FACT (front): %v\n", trimmed)
					fact = &Fact{Front: trimmed}
					facts = append(facts, fact)
				} else {
					if fact != nil {
						debugPlain("FACT (back): %v\n", trimmed)
						fact.Back = trimmed
					}
				}
			}

			if inHeader {
				debugPlain("HEADER: %v\n", trimmed)
				header = trimmed
			}

			buf.Reset()
			inHeader = false
			inNote = false
			indentCount = 0
			whitespaceCount = 0

		case "#":
			debug("#")

			if !inHeader && !inNote {
				inHeader = true
			} else {
				buf.WriteString(s.TokenText())
			}

		case "*":
			debug("*")

			if !inHeader && !inNote {
				inNote = true
			} else {
				buf.WriteString(s.TokenText())
			}

		default:
			debug("%s", s.TokenText())

			if inHeader || inNote {
				buf.WriteString(s.TokenText())
			}

			// any non-whitespace character resets whitespace count to zero
			whitespaceCount = 0
		}
	}

	return header, facts, nil
}

func main() {
	files, err := dataFiles()
	if err != nil {
		panic(err)
	}

	var facts []*Fact

	for _, file := range files {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			panic(err)
		}

		header, fileFacts, err := parseMarkdown(string(data))
		if err != nil {
			panic(err)
		}

		// add header as a tag to all facts
		if header != "" {
			for _, fact := range facts {
				fact.Tags = []string{header}
			}
		}

		facts = append(facts, fileFacts...)
	}

	var buf bytes.Buffer
	for _, fact := range facts {
		//fmt.Printf("fact: %v\n", *fact)
		buf.WriteString(fact.Front)
		buf.WriteString("\t")
		buf.WriteString(fact.Back)
		buf.WriteString("\t")
		buf.WriteString(strings.Join(fact.Tags, ","))
		buf.WriteString("\n")
	}

	fmt.Printf("writing: facts.tsv\n")
	err = ioutil.WriteFile("./facts.tsv", buf.Bytes(), 0644)
	if err != nil {
		panic(err)
	}
}
