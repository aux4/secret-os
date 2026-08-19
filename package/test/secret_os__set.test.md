# secret os set

## updating a field

```beforeAll
aux4 secret os create --vault Aux4Test --item Jira --fields "token=original,user=john" --index .secret/index.json > /dev/null 2>&1
```

```afterAll
aux4 secret os remove --ref Aux4Test/Jira --index .secret/index.json > /dev/null 2>&1 || true
rm -rf .secret
```

### should replace the value

```execute
aux4 secret os set --ref Aux4Test/Jira --field token --value rotated --index .secret/index.json
aux4 secret os get --ref Aux4Test/Jira --fields token --index .secret/index.json
```

```expect
secret://os/Aux4Test/Jira/token updated
{
  "token": "rotated"
}
```

### should leave the other fields untouched

```execute
aux4 secret os set --ref Aux4Test/Jira --field token --value rotated --index .secret/index.json > /dev/null
aux4 secret os get --ref Aux4Test/Jira --fields user --index .secret/index.json
```

```expect:json
{
  "user": "john"
}
```

### should refuse to create the item

```execute
aux4 secret os set --ref Aux4Test/Missing --field token --value v --index .secret/index.json
```

```error:partial
no secret found at secret://os/Aux4Test/Missing, use create first
```
