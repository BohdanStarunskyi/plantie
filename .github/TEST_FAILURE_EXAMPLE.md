# Testing CI Failure Protection

This file contains an example of how to test that the CI pipeline properly prevents merging when tests fail.

## Example failing test

To test that branch protection works with the CI pipeline, you can temporarily add a failing test like this:

```go
// In controllers/example_test.go (create this file temporarily)
package controllers

import "testing"

func TestExampleFailure(t *testing.T) {
    // This test will always fail to demonstrate CI protection
    t.Error("This is an intentional test failure to verify CI protection works")
}
```

## What should happen:

1. When you create a PR with this failing test, the CI will run
2. The `test` job will fail because of the failing test
3. If branch protection is configured correctly, GitHub will prevent merging the PR
4. The PR will show a red X indicating tests failed
5. You'll see a message like "All checks have failed" or "1 failing check"

## To fix:

Simply remove the failing test or fix it, and the CI will pass, allowing the PR to be merged.

## Note:

Remove this file and any example failing tests before finalizing your setup.