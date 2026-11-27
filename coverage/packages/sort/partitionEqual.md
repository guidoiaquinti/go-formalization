# Independent review of `sort.partitionEqual`

**Verdict: FAITHFUL**

Evidence SHA-256: `b8e8e6dce32f8543869628b2b130074dcf390d45ebb411e4f9578b6c9173e63c`

I traced the exact pinned Go 1.27.0 declaration alongside the Lean model and compared the test hypotheses with the source. The source expects a minimum-valued pivot. It moves equal elements ahead of greater elements via inward scans and returns the first greater index. An all-equal range returns the exclusive end while preserving the interior order.

I ran `go run ./difftest/sort/partition/review` on darwin/arm64. Ten independent literal traces matched native Go and compiled Lean. The witnesses reject an always-true partition flag, an omitted singleton self-swap, and an off-by-one equal boundary across both callback variants. The author compared 6,272 fixed rows and 60,000 samples at seeds 17, 2718 and 65537 for this declaration. [Reviewer control](../../../difftest/sort/partition/review/main.go)

The tested domain is a nonempty in-bounds range and a valid pivot. Less and Swap operate on signed integer array elements and return normally. A minimum-valued pivot is required for partitionEqual.
