# Independent review of `sort.stable`

**Verdict: FAITHFUL**

Evidence SHA-256: `84e9dd8441763fa1fdef4d1b1c2d5699ce48f374bad370401e1c0c2c5d310a87`

I traced the pinned Go 1.27.0 declaration, Lean model and native comparison. The pinned private entry sorts blocks of twenty, then doubles the merge width until the requested prefix is merged. The 21-element fixture distinguishes block merge from whole-prefix insertion.

I ran `go run ./difftest/sort/stable/review` on darwin/arm64. Ten independent observations matched compiled Lean and pinned Go. The checks distinguish reordered equal records, repeated Len, eager nil Less access and insertion across the whole 21-element range. The author replay covered 201 fixed cases and 60,000 sampled cases at seeds 17, 2718 and 65537 for this declaration. An independent oracle checks final records by key and original identity. [Reviewer control](../../../difftest/sort/stable/review/main.go)

The input scope covers representable indices, signed integer keys, and indexed callbacks on the tested record arrays. General callback mutation and arbitrary reflected element types are outside scope.
