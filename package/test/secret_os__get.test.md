# secret os get

## reading fields

```beforeAll
aux4 secret os create --vault Aux4Test --item Reader --fields "token=s3cr3t" --index .secret/index.json > /dev/null 2>&1
```

```afterAll
aux4 secret os remove --ref Aux4Test/Reader --index .secret/index.json > /dev/null 2>&1 || true
rm -rf .secret
```

### should return the field as json

```execute
aux4 secret os get --ref Aux4Test/Reader --fields token --index .secret/index.json
```

```expect:json
{
  "token": "s3cr3t"
}
```

### should fail on an unknown field

A mistyped reference has to fail loudly; resolving to nothing would hand the
caller an empty credential and fail somewhere far less obvious.

```execute
aux4 secret os get --ref Aux4Test/Reader --fields nosuchfield --index .secret/index.json
```

```error:partial
no secret found at secret://os/Aux4Test/Reader/nosuchfield
```

### should reject a reference without an item

```execute
aux4 secret os get --ref Aux4Test --fields token --index .secret/index.json
```

```error:partial
must be in the form <vault>/<item>
```
