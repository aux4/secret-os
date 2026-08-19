package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
)

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "Error: "+format+"\n", args...)
	os.Exit(1)
}

func argAt(args []string, index int) string {
	if len(args) > index {
		return strings.TrimSpace(args[index])
	}
	return ""
}

func splitList(raw string) []string {
	result := []string{}
	for _, part := range strings.Split(raw, ",") {
		part = strings.TrimSpace(part)
		if part != "" {
			result = append(result, part)
		}
	}
	return result
}

// splitRef splits a reference into vault and item: the vault is the first
// segment, the item is everything after it. This matches how the aux4 core
// splits a secret:// URI and how the other providers read the same reference.
func splitRef(ref string) (string, string, error) {
	ref = strings.Trim(strings.TrimSpace(ref), "/")
	if ref == "" {
		return "", "", fmt.Errorf("ref is required, in the form <vault>/<item>")
	}

	parts := strings.SplitN(ref, "/", 2)
	if len(parts) < 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("ref %q must be in the form <vault>/<item>", ref)
	}
	return parts[0], parts[1], nil
}

func parseFields(raw string) (map[string]string, error) {
	fields := map[string]string{}

	for _, pair := range strings.Split(raw, ",") {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}

		index := strings.Index(pair, "=")
		if index <= 0 {
			return nil, fmt.Errorf("field %q must be in the form name=value", pair)
		}

		name := strings.TrimSpace(pair[:index])
		if name == "" {
			return nil, fmt.Errorf("field %q has an empty name", pair)
		}
		fields[name] = pair[index+1:]
	}

	if len(fields) == 0 {
		return nil, fmt.Errorf("at least one field is required, in the form name=value")
	}
	return fields, nil
}

func reference(vault string, item string, field string) string {
	if field == "" {
		return fmt.Sprintf("secret://os/%s/%s", vault, item)
	}
	return fmt.Sprintf("secret://os/%s/%s/%s", vault, item, field)
}

// args: index ref fields otp
func runGet(args []string) {
	indexPath := argAt(args, 0)
	ref := argAt(args, 1)
	fields := splitList(argAt(args, 2))

	vault, item, err := splitRef(ref)
	if err != nil {
		fail("%v", err)
	}
	if len(fields) == 0 {
		fail("fields is required")
	}

	store := newKeystore()
	service := serviceName(vault, item)

	result := map[string]string{}
	names := []string{}

	for _, field := range fields {
		value, err := store.get(service, field)
		if err == errNotFound {
			fail("no secret found at %s in the %s", reference(vault, item, field), store.name())
		} else if err != nil {
			fail("%v", err)
		}
		result[field] = value
		names = append(names, field)
	}

	sort.Strings(names)
	lines := make([]string, 0, len(names))
	for _, name := range names {
		nameEncoded, _ := json.Marshal(name)
		valueEncoded, _ := json.Marshal(result[name])
		lines = append(lines, fmt.Sprintf("  %s: %s", nameEncoded, valueEncoded))
	}

	_ = indexPath
	fmt.Fprintf(os.Stdout, "{\n%s\n}\n", strings.Join(lines, ",\n"))
}

// args: index vault item fields category
func runCreate(args []string) {
	indexPath := argAt(args, 0)
	vault := argAt(args, 1)
	item := argAt(args, 2)
	rawFields := argAt(args, 3)

	if vault == "" {
		fail("vault is required")
	}
	if item == "" {
		fail("item is required")
	}
	if strings.Contains(vault, "/") {
		fail("vault %q must not contain '/', which separates the parts of a reference", vault)
	}

	fields, err := parseFields(rawFields)
	if err != nil {
		fail("%v", err)
	}
	for _, value := range fields {
		if err := validateValue(value); err != nil {
			fail("%v", err)
		}
	}

	index, err := loadIndex(indexPath)
	if err != nil {
		fail("%v", err)
	}
	if _, exists := index.fields(vault, item); exists {
		fail("%s already exists, use set to change a field", reference(vault, item, ""))
	}

	store := newKeystore()
	service := serviceName(vault, item)

	names := []string{}
	for name, value := range fields {
		if err := store.set(service, name, value); err != nil {
			fail("%v", err)
		}
		names = append(names, name)
	}

	index.add(vault, item, names)
	if err := index.save(indexPath); err != nil {
		fail("%v", err)
	}

	fmt.Fprintln(os.Stdout, reference(vault, item, ""))
}

// args: index ref field value
func runSet(args []string) {
	indexPath := argAt(args, 0)
	ref := argAt(args, 1)
	field := argAt(args, 2)
	value := ""
	if len(args) > 3 {
		value = args[3]
	}

	vault, item, err := splitRef(ref)
	if err != nil {
		fail("%v", err)
	}
	if field == "" {
		fail("field is required")
	}
	if err := validateValue(value); err != nil {
		fail("%v", err)
	}

	index, err := loadIndex(indexPath)
	if err != nil {
		fail("%v", err)
	}
	if _, exists := index.fields(vault, item); !exists {
		fail("no secret found at %s, use create first", reference(vault, item, ""))
	}

	store := newKeystore()
	if err := store.set(serviceName(vault, item), field, value); err != nil {
		fail("%v", err)
	}

	index.add(vault, item, []string{field})
	if err := index.save(indexPath); err != nil {
		fail("%v", err)
	}

	fmt.Fprintf(os.Stdout, "%s updated\n", reference(vault, item, field))
}

// args: index vault withFields
func runList(args []string) {
	indexPath := argAt(args, 0)
	vault := argAt(args, 1)
	withFields := argAt(args, 2) == "true"

	index, err := loadIndex(indexPath)
	if err != nil {
		fail("%v", err)
	}
	printEntries(index.entries(vault), withFields)
}

// args: index query vault withFields
func runSearch(args []string) {
	indexPath := argAt(args, 0)
	query := argAt(args, 1)
	vault := argAt(args, 2)
	withFields := argAt(args, 3) == "true"

	if query == "" {
		fail("query is required")
	}

	index, err := loadIndex(indexPath)
	if err != nil {
		fail("%v", err)
	}

	needle := strings.ToLower(query)
	matched := []indexEntry{}
	for _, entry := range index.entries(vault) {
		if strings.Contains(strings.ToLower(entry.Item), needle) {
			matched = append(matched, entry)
		}
	}
	printEntries(matched, withFields)
}

func printEntries(entries []indexEntry, withFields bool) {
	for _, entry := range entries {
		if !withFields {
			fmt.Fprintln(os.Stdout, reference(entry.Vault, entry.Item, ""))
			continue
		}
		for _, field := range entry.Fields {
			fmt.Fprintln(os.Stdout, reference(entry.Vault, entry.Item, field))
		}
	}
}

// args: index ref
func runRemove(args []string) {
	indexPath := argAt(args, 0)
	ref := argAt(args, 1)

	vault, item, err := splitRef(ref)
	if err != nil {
		fail("%v", err)
	}

	index, err := loadIndex(indexPath)
	if err != nil {
		fail("%v", err)
	}

	fields, exists := index.fields(vault, item)
	if !exists {
		fail("no secret found at %s", reference(vault, item, ""))
	}

	store := newKeystore()
	service := serviceName(vault, item)

	// Drop every field from the keystore before forgetting the item. Removing
	// the index entry first would strand the values with no way to name them.
	for _, field := range fields {
		if err := store.delete(service, field); err != nil && err != errNotFound {
			fail("%v", err)
		}
	}

	index.removeItem(vault, item)
	if err := index.save(indexPath); err != nil {
		fail("%v", err)
	}

	fmt.Fprintf(os.Stdout, "%s removed\n", reference(vault, item, ""))
}
