# Independent review of `bytes.Reader.Read`

**Verdict: FAITHFUL**

Evidence SHA-256: `c600a2759b4a37804c59dacadbc1a3ea53714d0a43926fabed679d9b8c72457b`

Reviewed author commit `19a3cae` against pinned Go 1.27.0 `src/bytes/reader.go:39–47`. End-of-input returns the exact `io.EOF` error before clearing the marker. A live read clears the marker, advances the cursor by the copied count, and returns nil even for a short copy. The copy uses the pre-write source bytes when the destination aliases the reader. The evidence compares the destination and full modeled backing storage, including spare bytes, as well as the cursor and marker.

The independent eight-case control passed in the pinned native overlay and Lean executable. It covers shared copies in both directions, a write into spare capacity, short disjoint copy with untouched destination suffix, empty live and EOF reads, nil receiver panic, and maximum cursor. Changing the Lean short-copy branch to return EOF failed two fixtures.

The claim covers valid slice geometry, nonnegative cursors, and sequential calls into disjoint or shared destination windows. Concurrency and invalid slice geometry remain outside the reviewed domain.
