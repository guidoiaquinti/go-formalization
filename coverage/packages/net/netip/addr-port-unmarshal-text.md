# Independent review of `net/netip.AddrPort.UnmarshalText`

**Verdict: FAITHFUL**

Reviewed author commit `752e268` against Go 1.27.0 `src/net/netip/netip.go`. Clears existing receiver state on ordinary parse failure and with successful empty input. Nonempty success stores the exact endpoint. The control distinguishes empty-success clearing from returning the parser error.

Pinned Go and Lean matched 574 fixed rows and 60,000 sampled rows using seeds 17, 2718 and 65537. Nine independent literal fixtures matched both. Eight altered outputs disagreed once each, and two separate Lean rows exercised 32- and 64-bit prefix conversion. Replay with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/net/netip/parse-endpoint/review` after building `godiff-netip-parse-endpoint`.

Failure comparison uses structured categories. Successful values compare family, bytes, raw zone and numeric field. Receiver checks use nonnil destinations.

Evidence SHA-256: `0aef2def0c3975112b4d28c394f3861ad43e67a6d6a3b9ef445d52c454a5295f`
