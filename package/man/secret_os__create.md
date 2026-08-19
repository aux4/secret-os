#### Description

The `create` command stores a new secret in the OS keystore and prints the `secret://` reference for it.

Each field becomes one keystore entry, keyed by `aux4:<vault>/<item>` with the field as the account name, so entries are identifiable in Keychain Access or Seahorse. The secret is passed to the platform tool on stdin, never as a command-line argument, so it does not appear in `ps`.

Values must be single-line: the macOS tool reads a password as a line, and a truncated credential is worse than a rejected one. Creating an item that already exists is an error; use `set` to change a field.


On Linux this needs a D-Bus session with an unlocked keyring, which containers, CI runners and plain SSH sessions do not have. On Windows the provider is not supported. In both cases use `aux4/secret-aux4`.

#### Usage

```bash
aux4 secret os create --vault <name> --item <title> --fields <key=value,...> [--category <type>] [--index <path>]
```

--vault      The vault namespace (required)
--item       The item title (required)
--fields     Comma-separated `key=value` assignments (required)
--category   The item category (default: `Login`)
--index      Path to the reference index (default: `${aux4HomeDir}/secret/os/index.json`)

#### Example

```bash
aux4 secret os create --vault Work --item "Billing API" --fields "clientId=my-client,clientSecret=s3cr3t"
```

```text
secret://os/Work/Billing API
```
