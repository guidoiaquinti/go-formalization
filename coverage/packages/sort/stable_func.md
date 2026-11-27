# Independent review of `sort.stable_func`

**Verdict: FAITHFUL**

Evidence SHA-256: `f8f9349da0ec51f9fc604685a2b52a71d0eda094e284c2657169621fec16726d`

I traced the pinned Go 1.27.0 declaration, Lean model and native comparison. The function callback variant uses the same twenty-element runs and merge schedule. The exact callback trace distinguishes equivalent final permutations achieved through another schedule.

I ran `go run ./difftest/sort/stable/review` on darwin/arm64. Ten independent observations matched compiled Lean and pinned Go. The checks distinguish reordered equal records, repeated Len, eager nil Less access and insertion across the whole 21-element range. The author replay covered 201 fixed cases and 60,000 sampled cases at seeds 17, 2718 and 65537 for this declaration. An independent oracle checks final records by key and original identity. [Reviewer control](../../../difftest/sort/stable/review/main.go)

The input scope covers representable indices, signed integer keys, and indexed callbacks on the tested record arrays. General callback mutation and arbitrary reflected element types are outside scope.
