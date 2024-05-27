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

| Verb              | `item` | `rel` | Effect                                                                  |
|-------------------|--------|-------|-------------------------------------------------------------------------|
| `create`          | X      | X     | Creates a new item or relationship.                                     |
| `fetch`           | X      | X     | Fetches an existing item or relationship.                               |
| `set`             | X      | X     | Sets the parameters of an existing item or relationship.                |
| `clear`           | X      | X     | Clears the parameters of an existing item or relationship.              |
| `delete`          | X      | X     | Deletes an existing item or relationship.                               |
| `list`            | X      | X     | Lists all items or relationships.                                       |
| `exists`          | X      | X     | Checks if an item or relationship exists.                               |
| `nest`            | X      |       | Nests an item within another item.                                      |
| `free`            | X      |       | Frees an item from its parent item.                                     |
| `in?`             | X      |       | Checks if an item is nested within another item.                        |
| `from?`           |        | X     | Checks if a relationship exists from a specified item.                  |
| `to?`             |        | X     | Checks if a relationship exists to a specified item.                    |
| `create-or-fetch` | X      | X     | Creates a new item or relationship, or fetches it if it already exists. |
| `create-or-set`   | X      | X     | Creates a new item or relationship, or sets if it already exists.       |

To regenerate the `pkg/grammar/grammar.peg.go` file:

```sh
go install github.com/pointlander/peg && \
peg -inline -switch -strict -output pkg/grammar/grammar.peg.go pkg/grammar/grammar.peg
```

#### Troubleshooting PEG

* `parse error near Whitespace (line 1 symbol 15 - line 1 symbol 16):`
    - If you've checked your grammar, and this should be a match, it's often that the match order isn't working for you.
    - See if something else could be matching your term before the error character that prevents the parser from going down the right path.
