/*
 * Copyright (c) 2022. The Greater Heptavirate (https://github.com/TheGreaterHeptavirate). All Rights Reservet
 *
 * All copies of this software (if not stated otherway) are dedicated
 * ONLY to personal, non-commercial use.
 */

package data

import (
	"errors"
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
func ReadTxt(fileData []byte) (*LegalAct, error) {
	var (
		result     = &LegalAct{}
		ruleText   = "" // *
		article    = "" // &
		subsection = "" // @
		chapter    = "" // %
		act        = "" // $

		isBlockEnd = true
	)

	add := func() {
		identifier := ""

		if chapter != "" {
			identifier += strings.Replace(chapter, "\n", ", ", len(chapter)-1)
		}

		if subsection != "" {
			identifier += strings.Replace(subsection, "\n", ", ", len(chapter)-1)
		}

		identifier += article

		result.Rules = append(result.Rules, Rule{
			Identifier: identifier,
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
				return nil, errors.New("more than one act in a file - not supported")
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
		case ruleMark:
			isBlockEnd = false
			ruleText += val + "\n"
		}
	}

	result.ActName = act

	return result, nil
}
