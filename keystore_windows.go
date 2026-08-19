package main

import "fmt"

// Windows Credential Manager is reachable only through the CredRead/CredWrite
// API -- cmdkey can store a generic credential but will not hand the value
// back, so there is no command-line path to a working implementation.
//
// Shipping the commands with an honest failure beats shipping untested P/Invoke:
// the user gets a clear reason and a working alternative instead of a silent
// misbehaviour. See the README for the aux4 provider, which works everywhere.
type windowsKeystore struct{}

func newKeystore() keystore { return windowsKeystore{} }

func (k windowsKeystore) name() string { return "Windows Credential Manager" }

func notSupported() error {
	return fmt.Errorf("the os provider does not support Windows yet; use the aux4 provider instead (aux4 aux4 pkger install aux4/secret-aux4)")
}

func (k windowsKeystore) get(service string, account string) (string, error) {
	return "", notSupported()
}

func (k windowsKeystore) set(service string, account string, value string) error {
	return notSupported()
}

func (k windowsKeystore) delete(service string, account string) error {
	return notSupported()
}
