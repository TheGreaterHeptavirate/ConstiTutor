/*
 * Copyright (c) 2022. The Greater Heptavirate (https://github.com/TheGreaterHeptavirate). All Rights Reserved
 *
 * All copies of this software (if not stated otherwise) are dedicated
 * ONLY to personal, non-commercial use.
 */

package app

import (
	"strings"

	"github.com/AllenDang/giu"
	"github.com/TheGreaterHeptavirate/ConstiTutor/pkg/data"
	"github.com/sahilm/fuzzy"
)

// Research performs a fuzzy search on the given phrase and updates the search results.
// if called for "", it will display all the rules in default order.
func (a *App) Research(phrase string) {
	// if the last character of phrase ar " " remove them (fuzzy search returns strange results)
	for strings.HasSuffix(phrase, " ") {
		phrase = phrase[:len(phrase)-1]
	}

	a.rows = make([]*giu.TableRowWidget, 0)

	src := make([]string, 0)

	for _, act := range a.data {
		for _, rule := range act.Rules {
			searchText := ""

			if a.searchOptions.actNames {
				searchText += act.ActName + " "
			}

			if a.searchOptions.paragraphs {
				searchText += rule.Identifier + " "
			}

			if a.searchOptions.text {
				searchText += rule.Text + " "
			}

			src = append(src, searchText)
		}
	}

	if phrase == "" {
		for _, act := range a.data {
			for _, rule := range act.Rules {
				if phrase == "" || strings.Contains(rule.Text, phrase) {
					a.addRow(rule)
				}
			}
		}
	}

	match := fuzzy.Find(phrase, src)
	for _, m := range match {
		_, rule := a.getRuleFromIndex(m.Index)
		a.addRow(*rule)
	}
}

func (a *App) getRuleFromIndex(i int) (actName string, rule *data.Rule) {
	for currentAct := 0; ; {
		if i < len(a.data[currentAct].Rules) {
			return a.data[currentAct].ActName, &a.data[currentAct].Rules[i]
		}

		i -= len(a.data[currentAct].Rules)
		currentAct++
	}
}

func (a *App) addRow(rule data.Rule) {
	a.rows = append(a.rows, giu.TableRow(
		giu.Label(rule.Identifier).Wrapped(true),
		giu.Label(rule.Text).Wrapped(true),
	))
}
