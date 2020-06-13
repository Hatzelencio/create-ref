## Create Ref

> Version: v0.0.1

------

## How to use it?

By default `create-ref` uses the env-variable `GITHUB_SHA` as a ref base to create a new reference. You only need set the variable `refs` with the value or values separates by commas. It only accepts a git refs format.

For more information you can check this [link](https://developer.github.com/v3/git/refs/#create-a-reference)

```yaml
jobs:
  job-id:
    runs-on: ubuntu-latest
    steps:
      - name: Create a ref
        uses: Hatzelencio/create-ref@v0.0.1
        with:
          refs: "tags/my-new-ref" # or refs: "heads/my-new-branch"
```

If you need specify the sha base, you can override the `sha` variable. Like the below sample:

```yaml
steps:
  - name: Create a ref
    uses: Hatzelencio/create-ref@v0.0.1
    with:
      refs: "heads/my-branch,tags/my-new-tag"
      sha: 8bbd7620d10bc2ac991db3d78cbcf2b868f76902
```

If you prefer, it's possible return an `exit code 1` if you set the variable `fail-if-ref-exists` with `FORCE`. By default, this variable is set with: `IGNORE`

```yaml
steps:
  - name: Create a ref
    uses: Hatzelencio/create-ref@v0.0.1
    with:
      refs: "heads/branch-red-already-exists"
      fail-if-ref-exists: FORCE
```