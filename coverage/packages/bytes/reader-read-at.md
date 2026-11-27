# Independent review of `bytes.Reader.ReadAt`

**Verdict: FAITHFUL**

Evidence SHA-256: `a6a8ea3cbaf4167cae192137bdd81310c52af9fd2f10d89b322cf3b8c6511b4f`

Reviewed author commit `19a3cae` against pinned Go 1.27.0 `src/bytes/reader.go:50–63`. A negative offset returns its fresh error before dereferencing the receiver. At and beyond end of input, even an empty destination receives the exact `io.EOF` error. A live call returns EOF on a short copy and preserves the reader cursor and rune marker. Writes through a shared destination remain visible in the reader storage.

The independent eight-case control passed in the pinned native overlay and Lean executable. It covers overlap in both directions, a short copy with preserved suffix, empty live and EOF read, maximum signed offset, and a nil receiver with either negative offset or nonnegative offset. A Lean mutant updating cursor and marker failed four fixtures. Checking receiver presence before validating the negative offset failed the nil negative-offset fixture.

The claim covers valid sequential slice geometries, nonnegative reader cursors, and signed offsets. Concurrency and invalid slice geometries remain outside the reviewed domain.
