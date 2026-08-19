# aux4/secret-os

Secret provider for `aux4/secret` backed by the operating system's own keystore — macOS Keychain, or the Secret Service on Linux (GNOME Keyring, KWallet).

The OS owns the key, so there is no key material to manage, and access is gated by the same login your desktop already protects.

## Installation

```bash
aux4 aux4 pkger install aux4/secret-os
```

## Quick Start

```bash
aux4 secret os create --vault Work --item "Billing API" --fields "clientId=my-client,clientSecret=s3cr3t"
# -> secret://os/Work/Billing API

aux4 curl request https://api.example.com --header "Authorization: Bearer secret://os/Work/Billing API/clientSecret"
```

You rarely call `get` yourself — put the `secret://` reference in a config file or an `.aux4` and aux4 resolves it when the command runs.

## Platform Support

| Platform | Backend | Status |
|----------|---------|--------|
| macOS | Keychain (`security`) | supported |
| Linux | Secret Service (`secret-tool`) | supported, requires a desktop session |
| Windows | Credential Manager | **not supported** |

**Windows** has no command-line path to read a stored credential back — `cmdkey` can write one but will not return the value. The commands are present and fail with a clear message rather than misbehaving.

**Linux** needs a D-Bus session with an unlocked keyring. That is normal on a desktop and absent in containers, on CI runners and over plain SSH; the provider says so explicitly rather than failing obscurely.

For those cases use [`aux4/secret-aux4`](https://hub.aux4.io/r/public/packages/aux4/secret-aux4), which works everywhere.

## How References Work

A reference is `secret://os/<vault>/<item>/<field>`.

The vault is a **namespace, not a container**. It is encoded into the key rather than mapped onto a real keychain or collection, so a reference committed to a shared config resolves on a machine that has never heard of that vault — which would be impossible if the vault had to exist first.

```text
secret://os/Work/Billing API/clientSecret
          │   │     │           └─ field
          │   │     └───────────── item
          │   └─────────────────── vault
          └─────────────────────── provider
```

Entries are stored under the key `aux4:<vault>/<item>` with the field as the account name, so they are identifiable in Keychain Access or Seahorse.

The item may contain `/`; the vault may not, since the first segment always separates them.

## Commands

### `aux4 secret os create`

```bash
aux4 secret os create --vault <name> --item <title> --fields <key=value,...> [--category <type>]
```

### `aux4 secret os get`

```bash
aux4 secret os get --ref <vault/item> --fields <field1,field2>
```

```json
{
  "clientId": "my-client",
  "clientSecret": "s3cr3t"
}
```

### `aux4 secret os set`

```bash
aux4 secret os set --ref <vault/item> --field <name> --value <value>
```

### `aux4 secret os list`

```bash
aux4 secret os list [--vault <name>] [--withFields <true|false>]
```

### `aux4 secret os search`

```bash
aux4 secret os search <query> [--vault <name>] [--withFields <true|false>]
```

### `aux4 secret os remove`

```bash
aux4 secret os remove --ref <vault/item>
```

## The Reference Index

Listing is not uniformly possible across these keystores — `security dump-keychain` prompts for permission item by item — so this package keeps its own record of which references exist, at `${aux4HomeDir}/secret/os/index.json`.

It holds **names only**: vault, item and field. No secret ever goes into it; the values stay in the keystore. Because the index is what `list` and `search` read, an entry deleted directly through Keychain Access still appears there until it is removed with `remove`. `get` always reads the keystore, so it remains authoritative.

## Limitations

- **Single-line values only.** The macOS tool reads a password as a line, so a value containing a line break is rejected rather than silently truncated. Use `aux4/secret-aux4` for multi-line secrets.
- **No one-time passwords.** `--otp` is accepted for contract compatibility but no TOTP seeds are stored.
- **Names are not secret.** Vault, item and field names are visible in the keystore and in the index.
