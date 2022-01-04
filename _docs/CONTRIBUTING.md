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

## Testing -fix

Tests inside `cmd/perfguard/testdata/quickfix` are structured like this:

* `cmd/perfguard/testdata/quickfix/$name` is a test directory.
* `cmd/perfguard/testdata/quickfix/$name/before.go` is a target Go file.
* `cmd/perfguard/testdata/quickfix/$name/out/after.go` is expected output file.

`before.go` file should be executable by `go run`, so it must be a main package.

QuickFix tests not only check for the replacements to be correct, it also runs
old and new code forms to check whether they produce identical results.

So, add some `println()` statements to your QuickFix tests to validate the semantics.
