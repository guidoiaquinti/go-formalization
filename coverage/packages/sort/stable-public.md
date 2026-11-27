# Independent review of `sort.Stable`

**Verdict: FAITHFUL**

Evidence SHA-256: `ad04fee9bf077cfa559d98ac23b04f30ac26c837d294cb84e3de1d2547d91a7f`

I traced the pinned Go 1.27.0 declaration, Lean model and native comparison. The public interface calls Len once and retains stable record identities within that prefix. Nil interfaces panic before any comparison.

I ran `go run ./difftest/sort/stable/review` on darwin/arm64. Ten independent observations matched compiled Lean and pinned Go. The checks distinguish reordered equal records, repeated Len, eager nil Less access and insertion across the whole 21-element range. The author replay covered 301 fixed cases and 60,000 sampled cases at seeds 17, 2718 and 65537 for this declaration. An independent oracle checks final records by key and original identity. [Reviewer control](../../../difftest/sort/stable/review/main.go)

The input scope covers representable indices, signed integer keys, and indexed callbacks on the tested record arrays. General callback mutation and arbitrary reflected element types are outside scope.
