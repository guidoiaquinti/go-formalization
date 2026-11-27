# Independent review of `bytes.Replace`

**Verdict: FAITHFUL**

Evidence SHA-256: `b7dbe1e67d001b53a88efebdf8c638630d318d5e2763d1321496b9796fb5b33f`

Reviewed author commit `01c772c` against pinned Go 1.27.0 `src/bytes/bytes.go:1162–1200`. Matches are nonoverlapping and stop at the signed limit when nonnegative. Empty patterns follow UTF-8 rune boundaries, with malformed bytes consuming one byte each. A no-work empty copy is nil. Removing every byte through an actual match yields a nonnil zero-capacity result. The native probe checks isolation on the zero-limit route with source-aliasing arguments.

Six literal fixtures passed against native Go and compiled Lean, alongside three native alias probes. Making every empty no-work copy nonnil changed two fixtures. Replay using `GOTOOLCHAIN=go1.27.0 go run ./difftest/bytes/copy-replace/review` from the repository root.

The author checked 25,653 boundaries and 60,000 sampled cases at seeds 17, 2718, and 65537. Exact capacity claims apply to zero-capacity results and replacement allocation. Allocator-dependent append capacity and allocation failure remain outside scope.
