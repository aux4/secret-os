# secret os list

## listing stored references

```beforeAll
aux4 secret os create --vault Aux4Test --item Alpha --fields "a=1" --index .secret/index.json > /dev/null 2>&1
aux4 secret os create --vault Aux4Test --item Beta --fields "b=2,c=3" --index .secret/index.json > /dev/null 2>&1
```

```afterAll
aux4 secret os remove --ref Aux4Test/Alpha --index .secret/index.json > /dev/null 2>&1 || true
aux4 secret os remove --ref Aux4Test/Beta --index .secret/index.json > /dev/null 2>&1 || true
rm -rf .secret
```

### should print references sorted by vault and item

```execute
aux4 secret os list --vault Aux4Test --index .secret/index.json
```

```expect
secret://os/Aux4Test/Alpha
secret://os/Aux4Test/Beta
```

### should print one line per field

```execute
aux4 secret os list --vault Aux4Test --withFields true --index .secret/index.json
```

```expect
secret://os/Aux4Test/Alpha/a
secret://os/Aux4Test/Beta/b
secret://os/Aux4Test/Beta/c
```

### should find an item by part of its title

```execute
aux4 secret os search alph --vault Aux4Test --index .secret/index.json
```

```expect
secret://os/Aux4Test/Alpha
```

### should print nothing for an empty index

```execute
aux4 secret os list --index .secret/empty.json
echo "(no output)"
```

```expect
(no output)
```
