# Independent review of `net/netip.uint128.addOne`

**Verdict: FAITHFUL**

Evidence SHA-256: `e900af5492d2c9311ac78c371f9159df3c52e91f7a3d72089d9aa79a7ba9acad`

I checked evidence commit `ab8244f` against Go 1.27.0 `net/netip/uint128.go`. An all-ones low limb wraps to zero and carries to the high limb. The independent fixture checks both the carry case and the full-word wrap.

I ran `go run ./difftest/net/netip/uint128` on darwin/arm64 and reproduced 1,796 boundaries and 20,000 cases per seed 17, 2718 and 65537. The compiled Lean and direct pinned private Go observations agreed with the arbitrary-precision reference. Five separate literal fixtures passed both engines with `go run ./difftest/net/netip/uint128/review`, with eleven fixed mutation checks across the group. The scope records panic classification and direct field-address behavior where applicable. Exact runtime fault text and concurrent unsafe access are outside scope. [Independent control](../../../../difftest/net/netip/uint128/review/main.go)
