#### Description

The `get` command reads fields from the OS keystore and prints them as a JSON object. This is the command the aux4 core invokes when it resolves a `secret://os/...` reference, so you rarely call it directly.

It reads the keystore, not the index, and is therefore authoritative: an entry removed directly through Keychain Access fails here even if it still appears in `list`.

A missing item or field is an error rather than an empty result, so a mistyped reference fails loudly instead of supplying an empty credential.


On Linux this needs a D-Bus session with an unlocked keyring, which containers, CI runners and plain SSH sessions do not have. On Windows the provider is not supported. In both cases use `aux4/secret-aux4`.

#### Usage

```bash
aux4 secret os get --ref <vault/item> --fields <field1,field2> [--otp <true|false>] [--index <path>]
```

--ref      The secret reference, without the provider or field (required)
--fields   Comma-separated field names (required)
--otp      Include a one-time password (not supported by this provider)
--index    Path to the reference index (default: `${aux4HomeDir}/secret/os/index.json`)

#### Example

```bash
aux4 secret os get --ref "Work/Billing API" --fields clientId,clientSecret
```

```json
{
  "clientId": "my-client",
  "clientSecret": "s3cr3t"
}
```
