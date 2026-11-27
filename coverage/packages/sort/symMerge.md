# Independent review of `sort.symMerge`

**Verdict: FAITHFUL**

Evidence SHA-256: `66bea6d838457a41c908a23e6e63d848adcaceb97549535466258748aa68f1e3`

I traced the pinned Go 1.27.0 declaration, Lean execution and same-package native overlay. The left singleton path searches for the first right value greater than or equal to its key. The right singleton path searches for the first strictly larger left key. The different tie tests retain the original order across the boundary. The general path computes symmetric cuts, rotates the middle blocks and recurses only into nondegenerate ranges.

I ran `go run ./difftest/sort/sym-merge/review` on darwin/arm64. Eight independent fixtures matched the exact native comparison and swap sequence and final element identities. Reversing either singleton tie rule or omitting the central rotate changed the observed traces, giving six variant-specific mutation witnesses. The author comparison passed 15,175 boundary fixtures and 60,000 samples at seeds 17, 2718 and 65537 for this declaration, with an independent linear stable merge oracle. [Reviewer control](../../../difftest/sort/sym-merge/review/main.go)

The domain is two sorted nonempty adjacent runs in signed integer storage, representable index arithmetic and ordinary indexed callbacks.
