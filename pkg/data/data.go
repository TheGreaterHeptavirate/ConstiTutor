// Package data contains logic used for importing json data about
// legal acts.
package data

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
)

// Load returns a list of legal act loaded from DATA directory.
func Load() ([]*LegalAct, error) {
	return loadDir(".")
}

func loadDir(dirname string) ([]*LegalAct, error) {
	result := make([]*LegalAct, 0)

	files, err := data.ReadDir(dirname)
	if err != nil {
		return nil, fmt.Errorf("reading directory %s: %w", dirname, err)
	}

	for _, file := range files {
		if file.IsDir() {
			dataFromDir, err := loadDir(joinPath(dirname, file.Name()))
			if err != nil {
				return nil, err
			}

			result = append(result, dataFromDir...)

			continue
		}

		filename := joinPath(dirname, file.Name())

		fileData, err := data.ReadFile(filename)
		if err != nil {
			return nil, fmt.Errorf("reading file %s: %w", filename, err)
		}

		output := &LegalAct{}

		switch filepath.Ext(file.Name()) {
		case ".json":
			err = json.Unmarshal(fileData, output)
			if err != nil {
				return nil, fmt.Errorf("unmarshaling JSON: %w", err)
			}
		case ".txt":
			output, err = ReadTxt(fileData)
		}

		result = append(result, output)
	}

	return result, nil
}

// LegalAct represents a set of Law rules (e.g. Konstytucja RP).
type LegalAct struct {
	// ActName is a name of the Legal act (e.g. Konstytucja RP)
	ActName string

	// Rules is a set of rules
	Rules []Rule
}

// Rule represents a single rule from a LegalAct.
type Rule struct {
	// Identifier is a "index" of article/paragraph
	// e.g.: "Rozdział 4, Artykuł 3, paragraf 7"
	Identifier string

	// Text is a text of the rule
	Text string

	// Links is a list to external resources. e.g. YouTube description e.t.c.
	Links []string
}

func joinPath(path ...string) string {
	return strings.ReplaceAll(filepath.Join(path...), "\\", "/")
}
