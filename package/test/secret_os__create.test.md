# secret os create

## storing a secret

```afterAll
aux4 secret os remove --ref "Aux4Test/Billing API" --index .secret/index.json > /dev/null 2>&1 || true
aux4 secret os remove --ref "Aux4Test/Once" --index .secret/index.json > /dev/null 2>&1 || true
aux4 secret os remove --ref "Aux4Test/Equals" --index .secret/index.json > /dev/null 2>&1 || true
rm -rf .secret
```

### should print the reference to the new secret

```execute
aux4 secret os create --vault Aux4Test --item "Billing API" --fields "clientId=my-client,clientSecret=s3cr3t" --index .secret/index.json
```

```expect
secret://os/Aux4Test/Billing API
```

### should read the values back from the keystore

```execute
aux4 secret os get --ref "Aux4Test/Billing API" --fields clientId,clientSecret --index .secret/index.json
```

```expect:json
{
  "clientId": "my-client",
  "clientSecret": "s3cr3t"
}
```

### should keep a value that contains an equals sign intact

```execute
aux4 secret os create --vault Aux4Test --item Equals --fields "token=abc=def==" --index .secret/index.json > /dev/null
aux4 secret os get --ref Aux4Test/Equals --fields token --index .secret/index.json
```

```expect:json
{
  "token": "abc=def=="
}
```

### should never write a secret into the index

The index records which references exist so that listing works the same on
every platform. The values stay in the keystore.

```execute
grep -c "s3cr3t" .secret/index.json || true
```

```expect
0
```

### should refuse to overwrite an existing item

```execute
aux4 secret os create --vault Aux4Test --item Once --fields "a=b" --index .secret/index.json > /dev/null
aux4 secret os create --vault Aux4Test --item Once --fields "a=c" --index .secret/index.json
```

```error:partial
secret://os/Aux4Test/Once already exists, use set to change a field
```

### should reject a vault containing a slash

```execute
aux4 secret os create --vault "Aux4Test/Sub" --item Thing --fields "a=b" --index .secret/index.json
```

```error:partial
must not contain '/'
```

### should reject a value containing a line break

A password is read as a single line, so a multi-line value would be silently
truncated into a credential that is quietly wrong.

```execute
aux4 secret os create --vault Aux4Test --item Multi --fields "token=line1
line2" --index .secret/index.json
```

```error:partial
must not contain a line break
```
