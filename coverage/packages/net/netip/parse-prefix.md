# Independent review of `net/netip.ParsePrefix`

**Verdict: FAITHFUL**

Reviewed author commit `752e268` against Go 1.27.0 `src/net/netip/netip.go`. Splits at the last slash, rejects zones and noncanonical prefix lengths, and retains host bits. The controls reject masking host bits and accepting leading zeros. An oversized decimal distinguishes 32-bit conversion error from 64-bit prefix bounds.

Pinned Go and Lean matched 574 fixed rows and 60,000 sampled rows using seeds 17, 2718 and 65537. Nine independent literal fixtures matched both. Eight altered outputs disagreed once each, and two separate Lean rows exercised 32- and 64-bit prefix conversion. Replay with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/net/netip/parse-endpoint/review` after building `godiff-netip-parse-endpoint`.

Failure comparison uses structured categories. Successful values compare family, bytes, raw zone and numeric field. Receiver checks use nonnil destinations.

Evidence SHA-256: `651d7badd7580cc64032fb2ad434e3e89c882a916ccb04392e20c1e3a903dffe`
