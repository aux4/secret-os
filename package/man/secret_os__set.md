#### Description

The `set` command replaces one field of an existing secret, leaving the others untouched. The value goes to the platform tool on stdin, never as an argument.

The item must already exist — creating one here would let a mistyped reference produce a half-populated item that fails much later, when something tries to resolve it.


On Linux this needs a D-Bus session with an unlocked keyring, which containers, CI runners and plain SSH sessions do not have. On Windows the provider is not supported. In both cases use `aux4/secret-aux4`.

#### Usage

```bash
aux4 secret os set --ref <vault/item> --field <name> --value <value> [--index <path>]
```

--ref     The secret reference (required)
--field   The field name to update (required)
--value   The new value, single-line (required)
--index   Path to the reference index (default: `${aux4HomeDir}/secret/os/index.json`)

#### Example

```bash
aux4 secret os set --ref "Work/Billing API" --field clientSecret --value new-s3cr3t
```

```text
secret://os/Work/Billing API/clientSecret updated
```
