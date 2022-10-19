package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/TheGreaterHeptavirate/ConstiTutor/pkg/data"
	"os"
	"strings"
)

const (
	commentMark    = '#'
	actMark        = '$'
	chapterMark    = '%'
	subsectionMark = '@'
	articleMark    = '&'
	ruleMark       = '*'
)

// reader implements io.Reader and provides ReadLine method
// which returns a line of text from the input file.
type reader struct {
	*strings.Reader
}

// ReadLine returns a line of text from the input file.
func (r *reader) ReadLine() (string, error) {
	var line string

	for {
		c, _, err := r.ReadRune()
		if err != nil {
			return "", err
		}

		if c == '\n' {
			break
		}

		line += string(c)
	}

	return line, nil
}

// THIS IS A SEPARATED SCRIPT NOT RELATED TO ANOTHER PACKAGES
// the script is intended to take @TomaszDyrka's txt file and convert
// it to json file(s) basing on data.LegalAct struct
//
// txt file layout:
// lines starting with # are ignored
// the following line starts are like header levels in HTML
// starting from the most important: $, %, &
// the lines starting with one of the above characters are considered
// titles of rules (particular rule name is a string-sum of all titles)
// the text of the rule starts with * character.
//
// The function prints out the json input of data.LegalAct
func main() {
	infile := flag.String("i", "konstytucjaRP.txt", "input file")
	outfile := flag.String("o", "konstytucjaRP.json", "output file")
	flag.Parse()

	fileData, err := os.ReadFile(*infile)
	if err != nil {
		panic(err)
	}

	var (
		result     = &data.LegalAct{}
		ruleText   = "" // *
		article    = "" // &
		subsection = "" // @
		chapter    = "" // %
		act        = "" // $

		isBlockEnd = true
	)

	add := func() {
		result.Rules = append(result.Rules, data.Rule{
			Identifier: fmt.Sprintf("%s,%s, %s", chapter, subsection, article),
			Text:       ruleText,
		})

		ruleText = ""
		isBlockEnd = true
	}

	r := &reader{strings.NewReader(string(fileData))}

	for line, err := r.ReadLine(); err == nil; line, err = r.ReadLine() {
		if len(line) == 0 {
			continue
		}

		firstChar := line[0]
		val := line[1:]

		switch firstChar {
		case commentMark:
			continue
		case actMark:
			if !isBlockEnd {
				panic("more than one act in a file - not supported")
			}

			act += val
		case chapterMark:
			if !isBlockEnd {
				add()

				chapter = ""
				subsection = ""
				article = ""
			}

			chapter += val + "\n"
		case subsectionMark:
			if !isBlockEnd {
				add()

				subsection = ""
				article = ""
			}

			subsection += val + "\n"
		case articleMark:
			if !isBlockEnd {
				add()

				article = ""
			}

			article += val + "\n"
		case '*':
			isBlockEnd = false
			ruleText += val + "\n"
		}
	}

	result.ActName = act

	output, err := json.MarshalIndent(result, "", "\t")
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(*outfile, output, 0644)
	if err != nil {
		panic(err)
	}
}
