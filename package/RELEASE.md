# aux4/secret-os 0.0.1

First release.

## Secrets in the operating system keystore

A `aux4/secret` provider backed by the platform's own keystore, so the OS owns the key and there is no key material to manage.

```bash
aux4 secret os create --vault Work --item "Billing API" --fields "clientId=my-client,clientSecret=s3cr3t"
# -> secret://os/Work/Billing API
```

| Platform | Backend | Status |
|----------|---------|--------|
| macOS | Keychain (`security`) | supported |
| Linux | Secret Service (`secret-tool`) | supported, needs a desktop session |
| Windows | Credential Manager | not supported |

Windows has no command-line path to read a stored credential back — `cmdkey` writes one but will not return the value. The commands are present and fail with a clear message rather than shipping untested behaviour.

On Linux, a missing D-Bus session (container, CI runner, plain SSH) is reported as such and points at `aux4/secret-aux4`, rather than surfacing the raw error.

## One path across three keystores

The three platforms address secrets differently — service plus account, arbitrary attributes, or a single flat target string. Windows is the tightest, so it sets the shape: the reference is joined into one opaque key, `aux4:<vault>/<item>`, with the field as the account name.

The vault is a **namespace, not a container**. It is encoded into the key rather than mapped onto a real keychain or collection, which is what lets a reference committed to a shared config resolve on a machine that has never heard of that vault.

## The reference index

Listing is not uniformly possible across these keystores — `security dump-keychain` prompts item by item, `cmdkey` will not return values — so the package records which references exist at `${aux4HomeDir}/secret/os/index.json`, holding **names only**. `get` always reads the keystore and stays authoritative.

## Notes

- Secrets are passed to the platform tools on stdin, never as arguments, so they do not appear in `ps`.
- Values must be single-line: a password is read as a line, and a truncated credential is worse than a rejected one.

## Per-platform packaging

One package name, one version — but a manifest per platform, so the Linux artifact declares `libsecret-tools` and the macOS one declares nothing. `secret://os/...` stays portable across a team while dependencies remain precise.
