# Independent review of `bytes.NewBufferString`

**Verdict: FAITHFUL**

Evidence SHA-256: `f9adca0cf6669f010b9b9734cf836decea04262d436648a69fc10fcb33c136af`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `b79078b4eee8639fcc3a8c248551cf2de47907b23c83627069281cf0643037c7`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `5a94d679e8aa0b3a2f75ef7523accd8a545ba8a09e03ab07e641cef8197f28ca`.

Go 1.27.0 constructs a Buffer with a string-to-byte-slice copy and leaves the other fields at zero. In the escaping case, the runtime obtains rounded capacity for small byte allocations and clears the spare storage. The empty conversion yields a nonnil zero-capacity slice. The Lean model preserves the raw bytes and represents the same initial cursor, read flag, nilness, and spare capacity.

The author compared 1,028 boundary inputs and 60,000 sampled inputs using seeds 17, 2718 and 65537. Nine independent pinned darwin/arm64 and Lean fixtures passed, including malformed bytes, nonnil empty output, and size-class transitions 8/9, 256/257, and 480/481. The control checks the returned view and confirms that mutating its first and last bytes leaves the input and a second constructor result unchanged. An exact-length capacity would fail four fixtures, and a nil empty view would fail one.

The model covers successful escaping allocations through 512 bytes. Other allocation sizes, escape-analysis outcomes, allocation failure, and native execution on other architectures remain outside this review.
