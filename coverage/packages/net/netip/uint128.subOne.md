# Independent review of `net/netip.uint128.subOne`

**Verdict: FAITHFUL**

Evidence SHA-256: `1ef3d67db1927a5930a7a82856ea7972cfc8bf3f5d4454fe0c10bc1852261118`

I checked evidence commit `ab8244f` against Go 1.27.0 `net/netip/uint128.go`. A zero low limb wraps to all ones and borrows from the high limb. The literal high-one low-zero case rejects subtracting only from the low limb.

I ran `go run ./difftest/net/netip/uint128` on darwin/arm64 and reproduced 1,796 boundaries and 20,000 cases per seed 17, 2718 and 65537. The compiled Lean and direct pinned private Go observations agreed with the arbitrary-precision reference. Five separate literal fixtures passed both engines with `go run ./difftest/net/netip/uint128/review`, with eleven fixed mutation checks across the group. The scope records panic classification and direct field-address behavior where applicable. Exact runtime fault text and concurrent unsafe access are outside scope. [Independent control](../../../../difftest/net/netip/uint128/review/main.go)
