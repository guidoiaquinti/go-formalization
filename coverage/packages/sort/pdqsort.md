# Independent review of `sort.pdqsort`

**Verdict: FAITHFUL**

Evidence SHA-256: `b2d751754d23571f69ebcdadde30cfb06474fb92913b74938e29f67c33bcbabf`

I traced the pinned Go 1.27.0 control flow alongside the Lean model and native overlay. The small-range branch precedes the zero-budget heap fallback. The model applies pattern breaking and decrements the budget only after imbalance, handles the descending hint before the already-sorted attempt, and follows the duplicate partition and recursive-side choices. Each step entering a child or continuing the loop excludes at least one element in the tested ordered domain, which justifies the model's range-length bound.

I ran `go run ./difftest/sort/pdqsort/review` on darwin/arm64. Ten fresh fixtures matched the exact native callback trace and Lean state. The independent predicates distinguish the two-element insertion trace, the zero-budget heap trace and positive-budget median trace at length thirteen, and the descending ninther and reversal trace at length fifty. Six variant-specific threshold or fallback witnesses rejected the wrong path. The author replay also passed 2,645 boundaries and 60,000 seeded cases at 17, 2718 and 65537 for this declaration. [Reviewer control](../../../difftest/sort/pdqsort/review/main.go)

Scope covers valid half-open ranges and indexed signed integer callbacks, a nonnegative initial budget, and an interior predecessor no greater than the selected elements. A universal sorting theorem, concurrent access and arbitrary callback effects are outside scope.
