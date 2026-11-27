# Independent review of `internal/bytealg.abigen_runtime_memequal`

**Verdict: FAITHFUL**

Evidence SHA-256: `cc96bc79fea85bd69f200165abc16286ae78c2fccb5b26704c5dc9f47f162ab4`

The pinned Go 1.27.0 declaration is linked to `runtime.memequal`. The overlay invokes that exact entry with readable initialized windows and a byte count, and the Lean model compares their corresponding prefixes. The observable result is true at size zero and otherwise reflects equality within the count. The contract concerns returned equality and leaves memory-access implementation details outside scope.

The author tested 5,125 boundary cases and 60,000 samples using seeds 17, 2718 and 65537. Nine independent native and Lean fixtures check nil and nonnil zero-size arguments, aliased initialized bytes, independent misalignments, and a mismatch just inside or outside the count. Comparing the complete provided views or dropping the last compared byte fails three cases each.

Pointer faults, invalid memory, timing and other unexecuted targets remain outside scope.
