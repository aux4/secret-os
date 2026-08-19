# secret os remove

## deleting a secret

```beforeAll
aux4 secret os create --vault Aux4Test --item Gone --fields "token=a,other=b" --index .secret/index.json > /dev/null 2>&1
```

```afterAll
aux4 secret os remove --ref Aux4Test/Gone --index .secret/index.json > /dev/null 2>&1 || true
rm -rf .secret
```

### should remove the reference and the stored values

```execute
aux4 secret os remove --ref Aux4Test/Gone --index .secret/index.json
aux4 secret os list --vault Aux4Test --index .secret/index.json
echo "(nothing left)"
```

```expect
secret://os/Aux4Test/Gone removed
(nothing left)
```

### should leave nothing readable behind

`get` reads the keystore rather than the index, so this proves the values are
really gone and not merely unlisted.

```execute
aux4 secret os remove --ref Aux4Test/Gone --index .secret/index.json > /dev/null
aux4 secret os get --ref Aux4Test/Gone --fields token --index .secret/index.json
```

```error:partial
no secret found at
```

### should fail on an unknown item

```execute
aux4 secret os remove --ref Aux4Test/Nope --index .secret/index.json
```

```error:partial
no secret found at secret://os/Aux4Test/Nope
```
