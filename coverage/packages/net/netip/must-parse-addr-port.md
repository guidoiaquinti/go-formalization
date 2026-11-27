# Independent review of `net/netip.MustParseAddrPort`

**Verdict: FAITHFUL**

Reviewed author commit `752e268` against Go 1.27.0 `src/net/netip/netip.go`. Returns the exact endpoint on success or panics with the corresponding error category. The empty-input control rejects returning an ordinary failure instead of panicking.

Pinned Go and Lean matched 574 fixed rows and 60,000 sampled rows using seeds 17, 2718 and 65537. Nine independent literal fixtures matched both. Eight altered outputs disagreed once each, and two separate Lean rows exercised 32- and 64-bit prefix conversion. Replay with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/net/netip/parse-endpoint/review` after building `godiff-netip-parse-endpoint`.

Failure comparison uses structured categories. Successful values compare family, bytes, raw zone and numeric field. Receiver checks use nonnil destinations.

Evidence SHA-256: `32023fa95d0fcd99d3d0205deadec9a21625b7c734e1589a97c6855db310d3ea`
