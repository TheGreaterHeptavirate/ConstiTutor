package app

import (
	"github.com/AllenDang/giu"
	"github.com/TheGreaterHeptavirate/ConstiTutor/pkg/data"
	"github.com/sahilm/fuzzy"
	"strings"
)

// Research performs a fuzzy search on the given phrase and updates the search results.
// if called for "", it will display all the rules in default order.
func (a *App) Research(phrase string) {
	// if the last character of phrase ar " " remove them (fuzzy search returns strange results)
	for strings.HasSuffix(phrase, " ") {
		phrase = phrase[:len(phrase)-1]
	}

	a.rows = make([]*giu.TableRowWidget, 0)
	src := search(a.data)

	if phrase == "" {
		for _, act := range a.data {
			for _, rule := range act.Rules {
				if phrase == "" || strings.Contains(rule.Text, phrase) {
					a.addRow(act.ActName, &rule)
				}
			}
		}
	}

	match := fuzzy.FindFrom(phrase, src)
	for _, m := range match {
		actName, rule := src.get(m.Index)
		a.addRow(actName, rule)
	}
}

type search []*data.LegalAct

func (s search) String(i int) string {
	for current := 0; ; {
		if i < len((s)[current].Rules) {
			return (s)[current].ActName + " " + s[current].Rules[i].Identifier + " " + s[current].Rules[i].Text
		}

		i -= len((s)[current].Rules)
		current++
	}
}

func (s search) get(i int) (actName string, rule *data.Rule) {
	for current := 0; ; {
		if i < len((s)[current].Rules) {
			return s[current].ActName, &s[current].Rules[i]
		}

		i -= len((s)[current].Rules)
		current++
	}
}

func (s search) Len() int {
	var l int
	for _, act := range s {
		l += len(act.Rules)
	}

	return l
}
