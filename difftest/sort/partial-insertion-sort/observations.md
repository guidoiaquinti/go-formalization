# Partial insertion with an offset

The comparison initially asserted that a true return meant the selected range
was sorted for every in-bounds invocation. The pinned private helper violates
that assertion when the input has an arbitrary prefix before the range.

Use `a = 2`, `b = 52`, and a length-54 input whose first five values are
`[100, 101, 0, 1, -1]` and whose remaining values are their indices.
The left shift reaches index zero because its inner loop is bounded by one,
independently of `a`. Both native variants return true and leave the first
five values as `[-1, 100, 101, 0, 1]`. The selected range starts with
`[101, 0, 1]`. The comparison requires this result and the exact callback trace.

The independent sortedness assertion applies to ranges starting at zero and
the short branch, where a left shift cannot cross a preceding prefix. The
differential comparison also covers arbitrary in-bounds offsets. It establishes
the helper's behavior on those inputs without making a claim about whether
such prefixes occur in the public sorting algorithm.
