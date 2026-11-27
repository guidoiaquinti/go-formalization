# Independent review of `net/netip.Prefix.UnmarshalText`

**Verdict: FAITHFUL**

Reviewed author commit `752e268` against Go 1.27.0 `src/net/netip/netip.go`. Clears existing state on parse failure and with successful empty input. A successful nonempty input stores the original address and prefix length. The control detects returning a parse failure on empty input.

Pinned Go and Lean matched 574 fixed rows and 60,000 sampled rows using seeds 17, 2718 and 65537. Nine independent literal fixtures matched both. Eight altered outputs disagreed once each, and two separate Lean rows exercised 32- and 64-bit prefix conversion. Replay with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/net/netip/parse-endpoint/review` after building `godiff-netip-parse-endpoint`.

Failure comparison uses structured categories. Successful values compare family, bytes, raw zone and numeric field. Receiver checks use nonnil destinations.

Evidence SHA-256: `5d3a305354966ccab97ac0460ba96b9a2928a9a030654879cd1e7abd65dbca75`
