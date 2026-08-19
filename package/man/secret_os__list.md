#### Description

The `list` command prints stored references, one per line, ready to paste into a config file.

It reads the reference index rather than the keystore. Enumerating the keystore directly is not uniformly possible — `security dump-keychain` prompts for permission item by item — so this package records which references exist. The index holds names only; no secret is ever written to it.

Because the index is a record rather than the keystore itself, an entry deleted directly through Keychain Access still appears here until `remove` is used. `get` reads the keystore and remains authoritative.


On Linux this needs a D-Bus session with an unlocked keyring, which containers, CI runners and plain SSH sessions do not have. On Windows the provider is not supported. In both cases use `aux4/secret-aux4`.

#### Usage

```bash
aux4 secret os list [--vault <name>] [--withFields <true|false>] [--index <path>]
```

--vault        Limit the listing to one vault (default: all vaults)
--withFields   Print one line per field instead of one per item (default: `false`)
--index        Path to the reference index (default: `${aux4HomeDir}/secret/os/index.json`)

#### Example

```bash
aux4 secret os list
```

```text
secret://os/Personal/GitHub
secret://os/Work/Billing API
```
