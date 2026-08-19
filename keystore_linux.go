package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type linuxKeystore struct{}

func newKeystore() keystore { return linuxKeystore{} }

func (k linuxKeystore) name() string { return "Secret Service" }

// attributes address an entry. secret-tool matches on the exact attribute set,
// so the same triple is used for every operation.
func attributes(service string, account string) []string {
	return []string{"service", service, "account", account}
}

// noSessionHint turns the least helpful failure mode into an actionable one.
// Secret Service needs a D-Bus session with an unlocked keyring, which is
// absent in containers, on CI runners and over plain SSH -- exactly where this
// will be tried and exactly where the raw error explains nothing.
func noSessionHint(stderr string) error {
	if strings.Contains(stderr, "Cannot autolaunch") ||
		strings.Contains(stderr, "DBUS_SESSION_BUS_ADDRESS") ||
		strings.Contains(stderr, "org.freedesktop.secrets") {
		return fmt.Errorf("no Secret Service session available (headless, container or CI). Use the aux4 provider instead: %s", strings.TrimSpace(stderr))
	}
	return fmt.Errorf("%s", strings.TrimSpace(stderr))
}

func (k linuxKeystore) get(service string, account string) (string, error) {
	command := exec.Command("secret-tool", append([]string{"lookup"}, attributes(service, account)...)...)

	var stdout, stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr

	if err := command.Run(); err != nil {
		// secret-tool exits non-zero with no output when nothing matches.
		if stderr.Len() == 0 {
			return "", errNotFound
		}
		return "", noSessionHint(stderr.String())
	}

	if stdout.Len() == 0 {
		return "", errNotFound
	}
	return strings.TrimRight(stdout.String(), "\n"), nil
}

func (k linuxKeystore) set(service string, account string, value string) error {
	args := append([]string{"store", "--label=" + service}, attributes(service, account)...)
	command := exec.Command("secret-tool", args...)
	// secret-tool reads the secret from stdin, keeping it off the argument list.
	command.Stdin = strings.NewReader(value)

	var stderr bytes.Buffer
	command.Stderr = &stderr

	if err := command.Run(); err != nil {
		return noSessionHint(stderr.String())
	}
	return nil
}

func (k linuxKeystore) delete(service string, account string) error {
	command := exec.Command("secret-tool", append([]string{"clear"}, attributes(service, account)...)...)

	var stderr bytes.Buffer
	command.Stderr = &stderr

	if err := command.Run(); err != nil {
		return noSessionHint(stderr.String())
	}
	return nil
}
