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
In your dedicated worktree, follow AGENTS.md and CONTRIBUTING.md to extend
coverage by one Go declaration. Run make inventory to inspect current coverage.
Keep UNSTARTED rows and sort them alphabetically by package import path, then
exact inventory symbol, using LC_ALL=C. Start with the first row and walk
inward past rows with existing evidence or declarations the Go driver cannot
call directly under the contributor guide. Coordinate before editing if you
converge on the same package or symbol.

Run make inventory TARGET=<selected-target> before editing. Add the model,
comparison and proof, then run make coverage-update and
make verify TARGET=<selected-target>. Leave passing new evidence at
AWAITING_REVIEW. Report the symbol, claim, domain or
case count and seeds, verification, current state and gaps. Link the evidence.
```

### Agent B

```text
In your dedicated worktree, follow AGENTS.md and CONTRIBUTING.md to extend
coverage by one Go declaration. Run make inventory to inspect current coverage.
Keep UNSTARTED rows and sort them alphabetically by package import path, then
exact inventory symbol, using LC_ALL=C. Start with the last row and walk
inward past rows with existing evidence or declarations the Go driver cannot
call directly under the contributor guide. Coordinate before editing if you
converge on the same package or symbol.

Run make inventory TARGET=<selected-target> before editing. Add the model,
comparison and proof, then run make coverage-update and
make verify TARGET=<selected-target>. Leave passing new evidence at
AWAITING_REVIEW. Report the symbol, claim, domain or
case count and seeds, verification, current state and gaps. Link the evidence.
```
