## Testing

When adding a new functionality, please make sure to add tests that cover it.

### Adding new rules to perfguard/_rules

`universal` and `lint` rules are easy to test.

Add a folder named `cmd/perfguard/testdata/rulestest/$name`, where `$name` is a name
of a new rule. Your new rule will be executed over all files inside that folder.

By convention, we name files inside that folder with the same name as the folder itself,
so the first file can be named `cmd/perfguard/testdata/rulestest/$name/$name.go`.

Now you can the tests with `go test -run /$name ./cmd/perfguard`.

If you want to run all rules, not just the new rule, omit the `-run` parameter.
