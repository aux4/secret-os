#### Description

The `os` secret provider stores secrets in the operating system's own keystore: the Keychain on macOS, or the Secret Service (GNOME Keyring, KWallet) on Linux. The OS owns the key, so there is no key material for you to manage.

A reference is `secret://os/<vault>/<item>/<field>`. The vault is a namespace encoded into the key rather than a real keychain or collection, so a reference committed to a shared config resolves on a machine that has never seen that vault.

Available subcommands: `create`, `get`, `set`, `list`, `search`, `remove`.


On Linux this needs a D-Bus session with an unlocked keyring, which containers, CI runners and plain SSH sessions do not have. On Windows the provider is not supported. In both cases use `aux4/secret-aux4`.

#### Usage

```bash
aux4 secret os <subcommand> [options]
```

#### Example

```bash
aux4 secret os create --vault Work --item "Billing API" --fields "clientId=my-client,clientSecret=s3cr3t"
```

```text
secret://os/Work/Billing API
```
