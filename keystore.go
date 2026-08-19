package main

import (
	"fmt"
	"strings"
)

// The OS keystores disagree about how a secret is addressed: macOS has a
// service plus an account, Linux has arbitrary attributes, Windows has a single
// flat target string. Windows is the tightest, so it sets the shape everyone
// else fits into -- a path joined into one opaque key, plus a field.
//
// Crucially the vault is NOT mapped onto a real keychain or collection. It is
// only a namespace inside the key. Mapping it to a physical container would
// demand that every teammate's machine already have a container by that name,
// which is exactly what stops a reference in a shared config from resolving.
const servicePrefix = "aux4:"

func serviceName(vault string, item string) string {
	return servicePrefix + vault + "/" + item
}

// keystore is the per-platform surface. Only these three operations touch the
// OS; listing is served from the index, so every platform behaves the same.
type keystore interface {
	get(service string, account string) (string, error)
	set(service string, account string, value string) error
	delete(service string, account string) error
	name() string
}

// errNotFound distinguishes "no such secret" from a keystore failure, so a
// missing entry can be reported precisely instead of as a generic error.
var errNotFound = fmt.Errorf("not found")

// validateValue rejects what the platform tools cannot round-trip.
//
// The macOS CLI reads a password as a line, so an embedded newline would be
// silently truncated -- storing a credential that is quietly wrong is far worse
// than refusing it.
func validateValue(value string) error {
	if strings.ContainsAny(value, "\n\r") {
		return fmt.Errorf("value must not contain a line break; the os provider stores single-line secrets, use the aux4 provider for anything else")
	}
	return nil
}
