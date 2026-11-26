# Extend coverage by one Go function

Copy the request below into a coding agent with this repository open. Leave either field as `auto` or specify a target and declaration. [Target list](../coverage/targets.json) · [Contributor guide](../CONTRIBUTING.md)

```text
Target = auto
Function = auto

Follow AGENTS.md and CONTRIBUTING.md to extend coverage by one Go declaration.
For Target = auto, use `make inventory` to find a configured target with an
eligible declaration. For a named target, search only that target. Before
editing, run `make inventory TARGET=<selected-target>`.

For Function = auto, choose an UNSTARTED inventory symbol and check for an
existing model first. For a named function or method, use its package import
path and exact inventory symbol, including the receiver for methods. If the
target is auto, select the target that contains that package and symbol.
Complete any missing evidence or link the current FAITHFUL review if the work
is already complete.

Report the symbol, claim, tested domain or seeds, verification results, final
inventory state and open gaps. Link the model, comparison and available review.
```

## Two agents working in parallel

Give each agent its own worktree from the same baseline commit. Each block below is a complete assignment. [Setup and integration](../CONTRIBUTING.md#parallel-coverage-work)

### Agent A

```text
Work in the dedicated coverage/utf8-runestart worktree. Follow AGENTS.md,
CONTRIBUTING.md and prompts/extend-coverage.md with Target = unicode/utf8
and Function = unicode/utf8.RuneStart. Run the target inventory before editing
and check for existing evidence. Compare every one of the 256 byte values
against Go and the claim that a rune can start at any byte outside 0x80–0xBF.
Run make coverage-update and make verify. Leave passing new evidence at
AWAITING_REVIEW. Report paths, case count, verification, inventory state and gaps.
```

### Agent B

```text
Work in the dedicated coverage/utf16-issurrogate worktree. Follow AGENTS.md,
CONTRIBUTING.md and prompts/extend-coverage.md with Target = unicode/utf16
and Function = unicode/utf16.IsSurrogate. Run the target inventory before editing
and check for existing evidence. Model the full signed rune domain and check
membership in U+D800 through U+DFFF inclusive. Include both bounds, adjacent
values, negative runes and int32 extremes, plus at least 50,000 reproducible
cases across several seeds. Run make coverage-update and make verify. Leave
passing new evidence at AWAITING_REVIEW. Report paths, counts, seeds,
verification, inventory state and gaps.
```
