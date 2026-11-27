# Independent review of `bytes.Repeat`

**Verdict: FAITHFUL**

Evidence SHA-256: `01617cc7e9789f584599e9c0ae73e49f3c2ea54374380f916f6defbcf42d14a3`

Reviewed author commit `5efb474` against pinned Go 1.27.0 `src/bytes/bytes.go:629–675`. The guards preserve the priority of a negative count on empty input. Products exceeding the signed maximum yield the exact overflow panic whether the product exceeds the unsigned word or fits it. For successful calls, the copy loop maintains source alignment after its 8 KiB transition and clips the result capacity to length.

Twelve independent fixtures passed against native Go and compiled Lean. The reviewer checks every successful byte at its original position modulo source length, nonnil empty results, exact capacity, and source and spare backing isolation. An early-empty mutant missed the negative-count panic. Treating overflow as a negative-count error changed two cases. Using an unaligned 8 KiB chunk changed both the large result of a three-byte source and the 8,193-byte source case. Replay using `GOTOOLCHAIN=go1.27.0 go run ./difftest/bytes/repeat-title/review` from the repository root.

Author evidence passed 638 boundaries and 60,000 sampled cases at seeds 17, 2718, and 65537. Runtime allocation failure after a representable product remains outside scope.
