# Testing fmt things

There was a bug in rpl that would perform the replacement work and if
`--show-diff` was specified, any fmt formatting string replacement
verbs would be erroneously rendered as "missing".

An example:

```
fmt.Printf("%v", `this should not be output as "missing"`)
```

Yet, after replacing text (correctly), the --show-diff output would
display `%!v(MISSING)` in the above statement.
