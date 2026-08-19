#### Description

The `search` command finds stored secrets whose item title contains the query, case-insensitively, and prints them as references. Only titles are searched — values live in the keystore and are never scanned.


On Linux this needs a D-Bus session with an unlocked keyring, which containers, CI runners and plain SSH sessions do not have. On Windows the provider is not supported. In both cases use `aux4/secret-aux4`.

#### Usage

```bash
aux4 secret os search <query> [--vault <name>] [--withFields <true|false>] [--index <path>]
```

query          Text to look for in item titles (required)
--vault        Limit the search to one vault (default: all vaults)
--withFields   Print one line per field instead of one per item (default: `false`)
--index        Path to the reference index (default: `${aux4HomeDir}/secret/os/index.json`)

#### Example

```bash
aux4 secret os search billing
```

```text
secret://os/Work/Billing API
```
