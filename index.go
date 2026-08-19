package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// The index records which references exist, and nothing else.
//
// It exists because listing is not uniformly possible across the platforms:
// `security dump-keychain` prompts for permission item by item, and `cmdkey`
// will not return values at all. Keeping our own record of the names makes
// list and search behave identically everywhere.
//
// It holds names only -- vault, item and field -- never a secret. Those names
// are not confidential in any of the providers, which is what allows a store to
// be browsed without unlocking it.
type Index struct {
	Version int                            `json:"version"`
	Vaults  map[string]map[string][]string `json:"vaults"`
}

func newIndex() *Index {
	return &Index{Version: 1, Vaults: map[string]map[string][]string{}}
}

func loadIndex(path string) (*Index, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return newIndex(), nil
		}
		return nil, fmt.Errorf("could not read the index: %v", err)
	}

	index := newIndex()
	if err := json.Unmarshal(data, index); err != nil {
		return nil, fmt.Errorf("the index at %s is not valid JSON: %v", path, err)
	}
	if index.Vaults == nil {
		index.Vaults = map[string]map[string][]string{}
	}
	return index, nil
}

func (i *Index) save(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}

	data, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		return err
	}

	temporary := path + ".tmp"
	if err := os.WriteFile(temporary, data, 0600); err != nil {
		return err
	}
	if err := os.Rename(temporary, path); err != nil {
		os.Remove(temporary)
		return err
	}
	return nil
}

func (i *Index) add(vault string, item string, fields []string) {
	if i.Vaults[vault] == nil {
		i.Vaults[vault] = map[string][]string{}
	}

	existing := map[string]bool{}
	for _, field := range i.Vaults[vault][item] {
		existing[field] = true
	}
	for _, field := range fields {
		if !existing[field] {
			i.Vaults[vault][item] = append(i.Vaults[vault][item], field)
			existing[field] = true
		}
	}
	sort.Strings(i.Vaults[vault][item])
}

func (i *Index) fields(vault string, item string) ([]string, bool) {
	items, ok := i.Vaults[vault]
	if !ok {
		return nil, false
	}
	fields, ok := items[item]
	return fields, ok
}

func (i *Index) removeItem(vault string, item string) bool {
	items, ok := i.Vaults[vault]
	if !ok {
		return false
	}
	if _, ok := items[item]; !ok {
		return false
	}

	delete(items, item)
	if len(items) == 0 {
		delete(i.Vaults, vault)
	}
	return true
}

type indexEntry struct {
	Vault  string
	Item   string
	Fields []string
}

func (i *Index) entries(vault string) []indexEntry {
	result := []indexEntry{}

	vaults := make([]string, 0, len(i.Vaults))
	for name := range i.Vaults {
		if vault != "" && name != vault {
			continue
		}
		vaults = append(vaults, name)
	}
	sort.Strings(vaults)

	for _, vaultName := range vaults {
		items := make([]string, 0, len(i.Vaults[vaultName]))
		for name := range i.Vaults[vaultName] {
			items = append(items, name)
		}
		sort.Strings(items)

		for _, itemName := range items {
			result = append(result, indexEntry{
				Vault:  vaultName,
				Item:   itemName,
				Fields: i.Vaults[vaultName][itemName],
			})
		}
	}

	return result
}
