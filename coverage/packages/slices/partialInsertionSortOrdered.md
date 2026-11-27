# Independent review of `slices.partialInsertionSortOrdered`

**Verdict: FAITHFUL**

Evidence SHA-256: `8d23f824ab3b59ce4d3127b79cc64e9a3c43fd82a173d8499fa2bd7dfa2a2557`

I checked evidence commit `d3db135` against the private Go 1.27.0 helper. It searches from `a + 1` and returns true when that scan reaches `b`. An empty interval returns false and a singleton returns true without consulting the comparator. On finding an inversion, only lengths at least fifty can repair it. A repair swaps the pair and shifts left and right. The left loop runs through index one, so it can modify index zero even with a positive range start. Reaching the five-repair limit returns false. Comparator failure preserves the observed callback prefix and completed swaps.

My `go run ./difftest/slices/partial-insertion` replay passed 25,742 boundaries and 20,000 cases each from seeds 17, 2718 and 65537 against pinned native Go and compiled Lean. Six independent literal fixtures passed with `go run ./difftest/slices/partial-insertion/review`. They reject shifting a length-49 inversion, clamping the left shift at the range start, rolling back an exchange on panic, and allowing a sixth repair. The prefix fixture uses start index two and moves the smallest value all the way to index zero on success. The panic while shifting left leaves that value at index one and preserves the comparison leading to the panic. [Independent control](../../../difftest/slices/partial-insertion/review/main.go)

The scope covers valid int64 ranges and fixed, non-mutating callback scripts. The boolean is the helper's own return flag, including executions with unconstrained prefixes. Arbitrary callback mutation, invalid bounds and concurrency are outside scope.
