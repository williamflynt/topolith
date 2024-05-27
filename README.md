# `topolith`

Interactive, expandable diagrams for complex systems.

#### Contents

...

---

...

## Code Structure

...

## Code Generation

...

### Grammar Generation

This package uses a PEG grammar, integrated with Go using the [peg](https://github.com/pointlander/peg) package by Andrew Snodgrass.
Thanks, Andrew!

To regenerate the `pkg/grammar/grammar.peg.go` file:

```sh
go install github.com/pointlander/peg && \
peg -inline -switch -strict -output pkg/grammar/grammar.peg.go pkg/grammar/grammar.peg
```

#### Troubleshooting PEG

* `parse error near Whitespace (line 1 symbol 15 - line 1 symbol 16):`
    - If you've checked your grammar, and this should be a match, it's often that the match order isn't working for you.
    - See if something else could be matching your term before the error character that prevents the parser from going down the right path.
