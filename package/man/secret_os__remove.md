#### Description

The `remove` command deletes every field of a secret from the OS keystore and then forgets the reference.

The keystore entries go first: dropping the index entry first would strand the values with no name left to address them by.

Removing a secret does not revoke it wherever it is used — it only forgets it locally.


On Linux this needs a D-Bus session with an unlocked keyring, which containers, CI runners and plain SSH sessions do not have. On Windows the provider is not supported. In both cases use `aux4/secret-aux4`.

#### Usage

```bash
aux4 secret os remove --ref <vault/item> [--index <path>]
```

--ref     The secret reference (required)
--index   Path to the reference index (default: `${aux4HomeDir}/secret/os/index.json`)

#### Example

```bash
aux4 secret os remove --ref "Work/Billing API"
```

```text
secret://os/Work/Billing API removed
```
