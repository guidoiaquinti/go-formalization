# Independent review of `net/netip.ParseAddrPort`

**Verdict: FAITHFUL**

Reviewed author commit `752e268` against Go 1.27.0 `src/net/netip/netip.go`. Splits at the last colon, validates the port before the address, then checks that the address family agrees with brackets. Leading decimal zeros are allowed for ports. The controls reject reversing fault precedence and accepting bracketed IPv4.

Pinned Go and Lean matched 574 fixed rows and 60,000 sampled rows using seeds 17, 2718 and 65537. Nine independent literal fixtures matched both. Eight altered outputs disagreed once each, and two separate Lean rows exercised 32- and 64-bit prefix conversion. Replay with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/net/netip/parse-endpoint/review` after building `godiff-netip-parse-endpoint`.

Failure comparison uses structured categories. Successful values compare family, bytes, raw zone and numeric field. Receiver checks use nonnil destinations.

Evidence SHA-256: `dae2fa4378d0f98f28e6bc67acaf00073dbc2cf741f651a1fa2fff98b3eaaa64`
