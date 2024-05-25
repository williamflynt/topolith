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

To regenerate the `pkg/grammar/repl.peg.go` file:

```sh
go install github.com/pointlander/peg && \
peg -inline -switch -strict -output pkg/grammar/grammar.peg.go pkg/grammar/grammar.peg
```
