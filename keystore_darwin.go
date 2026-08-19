package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type darwinKeystore struct{}

func newKeystore() keystore { return darwinKeystore{} }

func (k darwinKeystore) name() string { return "macOS Keychain" }

func (k darwinKeystore) get(service string, account string) (string, error) {
	command := exec.Command("security", "find-generic-password", "-s", service, "-a", account, "-w")

	var stdout, stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr

	if err := command.Run(); err != nil {
		if strings.Contains(stderr.String(), "could not be found") {
			return "", errNotFound
		}
		return "", fmt.Errorf("keychain read failed: %s", strings.TrimSpace(stderr.String()))
	}

	return strings.TrimRight(stdout.String(), "\n"), nil
}

func (k darwinKeystore) set(service string, account string, value string) error {
	// -w with no argument reads the password from the prompt, twice. Feeding it
	// on stdin keeps the secret out of the argument list, where `ps` would
	// expose it to every other user on the machine.
	command := exec.Command("security", "add-generic-password", "-U", "-s", service, "-a", account, "-w")
	command.Stdin = strings.NewReader(value + "\n" + value + "\n")

	var stderr bytes.Buffer
	command.Stderr = &stderr

	if err := command.Run(); err != nil {
		return fmt.Errorf("keychain write failed: %s", strings.TrimSpace(stderr.String()))
	}
	return nil
}

func (k darwinKeystore) delete(service string, account string) error {
	command := exec.Command("security", "delete-generic-password", "-s", service, "-a", account)

	var stderr bytes.Buffer
	command.Stderr = &stderr

	if err := command.Run(); err != nil {
		if strings.Contains(stderr.String(), "could not be found") {
			return errNotFound
		}
		return fmt.Errorf("keychain delete failed: %s", strings.TrimSpace(stderr.String()))
	}
	return nil
}
